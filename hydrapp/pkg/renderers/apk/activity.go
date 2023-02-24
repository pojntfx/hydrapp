package apk

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed activity.java
var activityTemplate string

type activityData struct {
	AppID string
}

func NewActivityRenderer(
	appID string,
) *renderers.Renderer {
	return renderers.NewRenderer("MainActivity.java", activityTemplate, activityData{appID})
}
