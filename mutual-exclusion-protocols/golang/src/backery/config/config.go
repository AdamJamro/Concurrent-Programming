package config

///////////////////////////////////////////
//
// define the parameters of the simulation
//
///////////////////////////////////////////

const NoOfProcesses = 10
const MinSteps = 30
const MaxSteps = 35

const MinDelay = 50 * millisecond  //nanoseconds
const MaxDelay = 100 * millisecond //nanoseconds

// the rest of configuration parameters are predefined
// and not meant to be altered

type ProcessState int

// this needs to be declared in order as the y coordinate
// uses iota values to indicate state
const (
	LocalSection ProcessState = iota
	EntryProtocol
	CriticalSection
	ExitProtocol // keep it as the last in the list
)

var SectionLabels = [NoOfSections]string{
	"LOCAL_SECTION",
	"ENTRY_PROTOCOL",
	"CRITICAL_SECTION",
	"EXIT_PROTOCOL",
}

const NoOfSections = int(ExitProtocol) + 1 // works under the assumption that ExitProtocol is the last listed state

const GridWidth = NoOfProcesses
const GridHeight = NoOfSections // works under the assumption that ExitProtocol is the last listed state

func ValidateConfiguration() {
	// TODO tailor for specific case
	if NoOfProcesses <= 0 {
		panic("NoOfProcesses must be greater than 0")
	}
	if MinSteps <= 0 {
		panic("MinSteps must be greater than 0")
	}
	if MaxSteps <= 0 {
		panic("MaxSteps must be greater than 0")
	}
	if MinSteps >= MaxSteps {
		panic("MinSteps must be strictly less than MaxSteps")
	}
}

const millisecond = 1000000
