package installer

import (
	"os"
	"regexp"

	"go.uber.org/zap"
)

var nilEncoderPreset = regexp.MustCompile(`(?m)^[ \t]*<EncoderPreset[^>]*xsi:nil="true"[^>]*/>[ \t]*\r?\n`)

// Jellyfin 10.11 wrote xsi:nil here; 12 made EncoderPreset a non-nullable enum, so the
// whole file fails to deserialize and the server saves defaults over the user's settings.
// Removable once no pre-12 device can still refresh in, which no release date guarantees.
func DropNilEncoderPreset(file string, logger *zap.Logger) error {
	before, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	after := nilEncoderPreset.ReplaceAll(before, nil)
	if len(after) == len(before) {
		return nil
	}

	logger.Info("dropping nil EncoderPreset", zap.String("file", file))

	info, err := os.Stat(file)
	if err != nil {
		return err
	}
	return os.WriteFile(file, after, info.Mode().Perm())
}
