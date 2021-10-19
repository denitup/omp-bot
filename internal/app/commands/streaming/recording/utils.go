package recording

import (
	"errors"
	"fmt"
	"strconv"
)

var errEmptyArgs = errors.New("empty command arguments")

func (c *RecordingCommander) getSingleUInt64Arg(strArg string) (uint64, error) {
	if strArg == "" {
		return 0, errEmptyArgs
	}

	intArg, err := strconv.ParseUint(strArg, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse uint64 argument: %w", err)
	}

	return intArg, nil
}

func (c *RecordingCommander) getRecordingIDArgument(arg string) (uint64, error) {
	recordingID, err := c.getSingleUInt64Arg(arg)
	if err == errEmptyArgs {
		return 0, errors.New("specify existing {recording_id} as an argument to the command")
	} else if err != nil {
		return 0, errors.New("invalid value for {recording_id} argument")
	}

	return recordingID, nil
}
