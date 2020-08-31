package kit

import (
	"github.com/ysmood/kit/pkg/http"
	"github.com/ysmood/kit/pkg/os"
	"github.com/ysmood/kit/pkg/run"
	"github.com/ysmood/kit/pkg/utils"
)

// All imported
var All = utils.All

// BackoffSleeper imported
var BackoffSleeper = utils.BackoffSleeper

// C imported
var C = utils.C

// ClearScreen imported
var ClearScreen = utils.ClearScreen

// CountSleeper imported
var CountSleeper = utils.CountSleeper

// DefaultBackoff imported
var DefaultBackoff = utils.DefaultBackoff

// Dump imported
var Dump = utils.Dump

// E imported
var E = utils.E

// E1 imported
var E1 = utils.E1

// Err imported
var Err = utils.Err

// ErrArg imported
var ErrArg = utils.ErrArg

// ErrInjector imported
type ErrInjector = utils.ErrInjector

// ErrMaxSleepCount imported
var ErrMaxSleepCount = utils.ErrMaxSleepCount

// JSON imported
var JSON = utils.JSON

// JSONResult imported
type JSONResult = utils.JSONResult

// Log imported
var Log = utils.Log

// MergeSleepers imported
var MergeSleepers = utils.MergeSleepers

// MustToJSON imported
var MustToJSON = utils.MustToJSON

// MustToJSONBytes imported
var MustToJSONBytes = utils.MustToJSONBytes

// Nil imported
type Nil = utils.Nil

// Noop imported
var Noop = utils.Noop

// Pause imported
var Pause = utils.Pause

// RandBytes imported
var RandBytes = utils.RandBytes

// RandString imported
var RandString = utils.RandString

// Retry imported
var Retry = utils.Retry

// S imported
var S = utils.S

// Sdump imported
var Sdump = utils.Sdump

// Sleep imported
var Sleep = utils.Sleep

// Sleeper imported
type Sleeper = utils.Sleeper

// Stderr imported
var Stderr = utils.Stderr

// Stdout imported
var Stdout = utils.Stdout

// Try imported
var Try = utils.Try

// Version imported
var Version = utils.Version

// GinContext imported
type GinContext = http.GinContext

// MustServer imported
var MustServer = http.MustServer

// Req imported
var Req = http.Req

// ReqContext imported
type ReqContext = http.ReqContext

// Server imported
var Server = http.Server

// ServerContext imported
type ServerContext = http.ServerContext

// CD imported
var CD = os.CD

// Chmod imported
var Chmod = os.Chmod

// Copy imported
var Copy = os.Copy

// DirExists imported
var DirExists = os.DirExists

// Escape imported
var Escape = os.Escape

// ExecutableExt imported
var ExecutableExt = os.ExecutableExt

// Exists imported
var Exists = os.Exists

// FileExists imported
var FileExists = os.FileExists

// HomeDir imported
var HomeDir = os.HomeDir

// Matcher imported
type Matcher = os.Matcher

// Mkdir imported
var Mkdir = os.Mkdir

// MkdirOptions imported
type MkdirOptions = os.MkdirOptions

// Move imported
var Move = os.Move

// NewMatcher imported
var NewMatcher = os.NewMatcher

// OutputFile imported
var OutputFile = os.OutputFile

// OutputFileOptions imported
type OutputFileOptions = os.OutputFileOptions

// ReadFile imported
var ReadFile = os.ReadFile

// ReadJSON imported
var ReadJSON = os.ReadJSON

// ReadString imported
var ReadString = os.ReadString

// Remove imported
var Remove = os.Remove

// RemoveWithDir imported
var RemoveWithDir = os.RemoveWithDir

// RetryPanic imported
var RetryPanic = os.RetryPanic

// SendSigInt imported
var SendSigInt = os.SendSigInt

// WaitSignal imported
var WaitSignal = os.WaitSignal

// Walk imported
var Walk = os.Walk

// WalkContext imported
type WalkContext = os.WalkContext

// WalkDirent imported
type WalkDirent = os.WalkDirent

// WalkFunc imported
type WalkFunc = os.WalkFunc

// WalkGitIgnore imported
var WalkGitIgnore = os.WalkGitIgnore

// WalkIgnoreHidden imported
var WalkIgnoreHidden = os.WalkIgnoreHidden

// Exec imported
var Exec = run.Exec

// ExecContext imported
type ExecContext = run.ExecContext

// GoBin imported
var GoBin = run.GoBin

// GoPath imported
var GoPath = run.GoPath

// Guard imported
var Guard = run.Guard

// GuardContext imported
type GuardContext = run.GuardContext

// GuardDefaultPatterns imported
var GuardDefaultPatterns = run.GuardDefaultPatterns

// KillTree imported
var KillTree = run.KillTree

// LookPath imported
var LookPath = run.LookPath

// MustGoTool imported
var MustGoTool = run.MustGoTool

// Task imported
var Task = run.Task

// TaskCmd imported
type TaskCmd = run.TaskCmd

// TaskContext imported
type TaskContext = run.TaskContext

// Tasks imported
var Tasks = run.Tasks

// TasksContext imported
type TasksContext = run.TasksContext

// TasksNew imported
var TasksNew = run.TasksNew
