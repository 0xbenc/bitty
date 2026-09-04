_bitty() {
  local cur="${COMP_WORDS[COMP_CWORD]}"
  COMPREPLY=($(compgen -W '--width --height --density --speed --seed --wrap --paused --intro --no-intro --no-color --theme-file --version --help' -- "$cur"))
}
complete -F _bitty bitty
