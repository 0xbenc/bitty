#compdef bitty
_arguments \
  '--width[world width]:cells:' '--height[world height]:cells:' \
  '--density[initial live-cell density]:density:' '--speed[generations per second]:speed:' \
  '--seed[random seed]:seed:' '--wrap[connect opposite edges]:boolean:(true false)' \
  '--paused[start paused]' '--intro[play startup intro]' '--no-intro[skip startup intro]' \
  '--no-color[disable color]' '--theme-file[portable theme file]:file:_files' \
  '--version[print version]' '--help[show help]'
