package display

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"syscall"

	"golang.org/x/exp/slices"
)

type mutex = sync.Mutex
type Chromium struct {
	mutex

	command *exec.Cmd
	input   *os.File
	output  *os.File

	lastMessage     uint32
	pendingMessages map[uint32]chan any

	targetId  string
	sessionId string
	windowId  uint
}

type jsonMap = map[string]any
type protoMsg struct {
	Id        uint32  `json:"id"`
	SessionId string  `json:"sessionId,omitempty"`
	Method    string  `json:"method,omitempty"`
	Params    jsonMap `json:"params,omitempty"`
	Result    jsonMap `json:"result,omitempty"`
	Error     jsonMap `json:"error,omitempty"`
}

func NewChromium(browser, url string, posX, posY, width, height uint) (*Chromium, error) {
	args := []string{
		fmt.Sprintf("--app=%s", url),
		fmt.Sprintf("--window-position=%d,%d", posX, posY),
		fmt.Sprintf("--window-size=%d,%d", width, height),
		"--remote-debugging-pipe",
		"--incognito",
		"--hide-scrollbars",
		"--deny-permission-prompts",
		"--block-new-web-contents",
		"--autoplay-policy=no-user-gesture-required",
		"--noerrdialogs",
		"--metrics-recording-only",
		"--password-store=basic",
		"--use-mock-keychain",
		"--no-first-run",
		"--no-default-browser-check",
		"--no-experiments",
		"--no-pings",
		"--disable-client-side-phishing-detection",
		"--disable-component-extensions-with-background-pages",
		"--disable-default-apps",
		"--disable-extensions",
		"--disable-external-intent-requests",
		"--disable-notifications",
		"--disable-popup-blocking",
		"--disable-prompt-on-repost",
		"--disable-background-networking",
		"--disable-breakpad",
		"--disable-domain-reliability",
		"--disable-sync",
		"--disable-back-forward-cache",
		"--disable-features=InterestFeedContentSuggestions",
		"--disable-features=Translate",
		"--disable-features=GlobalMediaControls",
		"--disable-features=ImprovedCookieControls",
		"--disable-features=MediaRouter",
		"--disable-features=OptimizationHints",
		"--disable-features=AutofillServerCommunication",
		"--disable-features=DialMediaRouteProvider",
		"--disable-features=BackForwardCache",
		"--disable-features=AvoidUnnecessaryBeforeUnloadCheckSync",
	}

	inputReader, inputWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	outputReader, outputWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(browser, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{inputReader, outputWriter}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Gid: 1000, Uid: 1000},
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Chromium{
		command: cmd,
		input:   inputWriter,
		output:  outputReader,

		pendingMessages: make(map[uint32]chan any),
	}, nil
}

func (chromium *Chromium) StartConnectionHandler() error {
	buffer := make([]byte, 1024)
	pending := make([]byte, 0)

	for {
		length, err := chromium.output.Read(buffer)
		if err != nil {
			return err
		}

		read := buffer[:length]
		start := 0
		for {
			end := slices.Index(read[start:], 0)
			if end == -1 {
				pending = append(pending, read[start:]...)
				break
			}

			end += start
			pending = append(pending, read[start:end]...)

			if err := chromium.receiveProtocolMessage(pending); err != nil {
				return err
			}

			start = end + 1
			pending = make([]byte, 0)
		}
	}
}

func (chromium *Chromium) InitConnection() error {
	targets, err := chromium.sendProtocolMessage("Target.getTargets", nil)
	if err != nil {
		return err
	}

	targetInfos, ok := targets["targetInfos"].([]any)
	if !ok {
		return errors.New("Target.getTargets: Unexpected targetInfos")
	}

	for _, targetInfo := range targetInfos {
		target, ok := targetInfo.(jsonMap)
		if !ok {
			return errors.New("Target.getTargets: Unexpected targetInfo")
		}

		t, ok := target["type"].(string)
		if !ok {
			return errors.New("Target.getTargets: Unexpected type")
		}

		if t != "page" {
			continue
		}

		targetId, ok := target["targetId"].(string)
		if !ok {
			return errors.New("Target.getTargets: Unexpected targetId")
		}

		chromium.targetId = targetId
		break
	}

	session, err := chromium.sendProtocolMessage("Target.attachToTarget", jsonMap{"targetId": chromium.targetId, "flatten": true})
	if err != nil {
		return err
	}

	sessionId, ok := session["sessionId"].(string)
	if !ok {
		return errors.New("Target.attachToTarget: Unexpected sessionId")
	}

	chromium.sessionId = sessionId

	_, err = chromium.sendProtocolMessage("Target.setAutoAttach", jsonMap{
		"autoAttach":             true,
		"flatten":                true,
		"waitForDebuggerOnStart": false,
	})
	if err != nil {
		return err
	}

	window, err := chromium.sendProtocolMessage("Browser.getWindowForTarget", jsonMap{"targetId": chromium.targetId})
	if err != nil {
		return err
	}

	windowId, ok := window["windowId"].(float64)
	if !ok {
		return errors.New("Browser.getWindowForTarget: Unexpected windowId")
	}

	chromium.windowId = uint(windowId)

	return nil
}

func (chromium *Chromium) Load(url string) error {
	_, err := chromium.sendProtocolMessage("Page.navigate", jsonMap{"url": url})
	return err
}

func (chromium *Chromium) SetPosition(posX, posY uint) error {
	_, err := chromium.sendProtocolMessage("Browser.setWindowBounds", jsonMap{
		"windowId": chromium.windowId,
		"bounds": jsonMap{
			"left": posX,
			"top":  posY,
		},
	})

	return err
}

func (chromium *Chromium) SetSize(width, height uint) error {
	_, err := chromium.sendProtocolMessage("Browser.setWindowBounds", jsonMap{
		"windowId": chromium.windowId,
		"bounds": jsonMap{
			"width":  width,
			"height": height,
		},
	})

	return err
}

func (chromium *Chromium) Kill() error {
	if chromium.command.ProcessState == nil {
		return chromium.command.Process.Kill()
	}

	return nil
}

func (chromium *Chromium) Wait() error {
	return chromium.command.Wait()
}

func (chromium *Chromium) sendProtocolMessage(method string, params jsonMap) (jsonMap, error) {
	id := atomic.AddUint32(&chromium.lastMessage, 1)

	msg, err := json.Marshal(protoMsg{
		Id:        id,
		SessionId: chromium.sessionId,
		Method:    method,
		Params:    params,
	})
	if err != nil {
		return nil, err
	}

	resChan := make(chan any)

	chromium.Lock()
	chromium.pendingMessages[id] = resChan
	chromium.Unlock()

	_, err = chromium.input.Write(append(msg, 0))
	if err != nil {
		return nil, err
	}

	result := <-resChan

	if res, ok := result.(jsonMap); ok {
		return res, nil
	}

	if err, ok := result.(error); ok {
		return nil, err
	}

	return nil, errors.New(method + ": Unexpected result")
}

func (chromium *Chromium) receiveProtocolMessage(messageBuffer []byte) error {
	msg := new(protoMsg)
	if err := json.Unmarshal(messageBuffer, msg); err != nil {
		return err
	}

	if msg.SessionId != chromium.sessionId {
		return nil
	}

	chromium.Lock()
	resChan, ok := chromium.pendingMessages[msg.Id]
	delete(chromium.pendingMessages, msg.Id)
	chromium.Unlock()

	if !ok {
		return nil
	}

	if msg.Result != nil {
		resChan <- msg.Result
		return nil
	}

	if msg.Error != nil {
		resChan <- msg.Result
		return nil
	}

	return errors.New("malformed protocol message")
}
