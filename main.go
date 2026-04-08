package main

import (
	_ "github.com/nuxtblog/nuxtblog/internal/packed"

	_ "github.com/nuxtblog/nuxtblog/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/nuxtblog/nuxtblog/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
