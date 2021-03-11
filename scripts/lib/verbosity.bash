#!/user/bib/env bash

VERBOSITY_CURRENT=0
VERBOSITY_FATAL=-3
VERBOSITY_ERROR=-2
VERBOSITY_QUIET=-1
VERBOSITY_ECHO=0
VERBOSITY_WARN=0
VERBOSITY_NOTICE=1
VERBOSITY_INFO=2
VERBOSITY_DEBUG=3
VERBOSITY_HELL=6

VERBOSITY_ERROR_FORMAT='\e[100\e[91m%s%s\e[0m: %s'
VERBOSITY_FATAL_FORMAT='\e[41m%s%s\e[0m: %s'
VERBOSITY_QUIET_FORMAT='\e[94m%s%s\e[0m: %s'
VERBOSITY_WARN_FORMAT='\e[91m%s%s\e[0m: %s'
VERBOSITY_NOTICE_FORMAT='\e[93m%s%s\e[0m: %s'
VERBOSITY_INFO_FORMAT='\e[96m%s%s\e[0m: %s'
VERBOSITY_DEBUG_FORMAT='\e[94m%s%s\e[0m: %s'
VERBOSITY_HELL_FORMAT='\e[102m\e[95m\e[5m\e[30m\e[2m'

function echo.verbosity.init_optparse() {
    VERBOSITY_CURRENT=0
    
    declare -F 'optparse.define' > /dev/null
    
    [[ $? -eq 0 ]] && {
        optparse.define short=v description="Verbosity" dispatch="echo.verbosity.raise"
        optparse.define short=q description="Quiet Mode" dispatch="echo.verbosity.quiet"
    } || {
        echo.warn "OptParse is not available, can not initialize verbosity handlers."    
        optparse.throw_error "um...  but I am..."
    }
}
    
function echo.verbosity.raise(){
    [[ $VERBOSITY_CURRENT -ge $VERBOSITY_HELL ]] && { echo.fatal "OK, THAT'S ENOUGH. I'M OUT"; } && return;
    [[ $VERBOSITY_CURRENT -lt 0 ]] && VERBOSITY_CURRENT="0"
    (( ++VERBOSITY_CURRENT ))
    [[ $VERBOSITY_CURRENT -eq $VERBOSITY_HELL ]] && {
        VERBOSITY_ERROR_FORMAT="$VERBOSITY_HELL_FORMAT$( echo $VERBOSITY_ERROR_FORMAT | sed -e 's/\\e\[0m//g' )"
        VERBOSITY_FATAL_FORMAT="$VERBOSITY_HELL_FORMAT$( echo $VERBOSITY_FATAL_FORMAT | sed -e 's/\\e\[0m//g' )"
        VERBOSITY_QUIET_FORMAT="$VERBOSITY_HELL_FORMAT$( echo $VERBOSITY_QUIET_FORMAT | sed -e 's/\\e\[0m//g' )"
        VERBOSITY_WARN_FORMAT="$VERBOSITY_HELL_FORMAT$( echo $VERBOSITY_WARN_FORMAT | sed -e 's/\\e\[0m//g' )"
        VERBOSITY_NOTICE_FORMAT="$VERBOSITY_HELL_FORMAT$( echo $VERBOSITY_NOTICE_FORMAT | sed -e 's/\\e\[0m//g' )"
        VERBOSITY_INFO_FORMAT="$VERBOSITY_HELL_FORMAT$( echo $VERBOSITY_INFO_FORMAT | sed -e 's/\\e\[0m//g' )"
        VERBOSITY_DEBUG_FORMAT="$VERBOSITY_HELL_FORMAT$( echo $VERBOSITY_DEBUG_FORMAT | sed -e 's/\\e\[0m//g' )"
        VERBOSITY_PREFIX="HELL MODE $VERBOSITY_PREFIX";
        echo.warn "Oh, right, hell mode...  SURPRISE...\n\tYou brought this on yourself...\n\tTHERE IS NO RETURN... THIS IS YOUR LIFE NOW...";
    } && return;
    [[ $VERBOSITY_CURRENT -gt $VERBOSITY_DEBUG ]] && echo.notice "We can't really get more verbose than this, sorry..." && return;
    [[ $VERBOSITY_CURRENT -eq $VERBOSITY_DEBUG ]] && echo.debug "You have entered DEBUG levels of verbosity... God Speed..." && return;
    true
}

function echo.verbosity.quiet(){
    [[ $VERBOSITY_CURRENT -ge $VERBOSITY_HELL ]] && echo.warn "I WON'T QUIET DOWN." && return;
    VERBOSITY_CURRENT=$VERBOSITY_QUIET
    #echo.quiet "We don't really support QUIET mode, but we'll try, just for you."
}

function echo(){
    local urgency=0
    [[ ! -z "$__VERBOSITY_URGENCY" ]] && urgency="$__VERBOSITY_URGENCY"
    unset __VERBOSITY_URGENCY
    [[ $VERBOSITY_CURRENT -lt $urgency ]] && return
    builtin echo "$@"
    [[ $urgency -le $VERBOSITY_FATAL ]] && { exit 1; } || true
}

function echo.quiet(){
    __VERBOSITY_URGENCY=$VERBOSITY_QUIET
    printf -v message "$VERBOSITY_QUIET_FORMAT\e[0m" "$VERBOSITY_PREFIX" "WHISPER" "$1"
    echo -e "$message"
}

function echo.error(){
    __VERBOSITY_URGENCY=$VERBOSITY_ERROR
    printf -v message "$VERBOSITY_ERROR_FORMAT\e[0m" "$VERBOSITY_PREFIX" "ERROR" "$1"
    >&2 echo -e "$message"
}

function echo.fatal(){
    __VERBOSITY_URGENCY=$VERBOSITY_FATAL
    printf -v message "$VERBOSITY_FATAL_FORMAT\e[0m" "$VERBOSITY_PREFIX" "FATAL" "$1"
    >&2 echo -e "$message"
}

function echo.warn(){
    __VERBOSITY_URGENCY=$VERBOSITY_WARN
    printf -v message "$VERBOSITY_WARN_FORMAT\e[0m" "$VERBOSITY_PREFIX" "WARN" "$1"
    echo -e "$message"
}

function echo.notice(){
    __VERBOSITY_URGENCY=$VERBOSITY_NOTICE
    printf -v message "$VERBOSITY_NOTICE_FORMAT\e[0m" "$VERBOSITY_PREFIX" "NOTICE" "$1"
    echo -e "$message"
}

function echo.info(){
    __VERBOSITY_URGENCY=$VERBOSITY_INFO
    printf -v message "$VERBOSITY_INFO_FORMAT\e[0m" "$VERBOSITY_PREFIX" "INFO" "$1"
    echo -e "$message"
}

function echo.debug(){
    __VERBOSITY_URGENCY=$VERBOSITY_DEBUG
    printf -v message "$VERBOSITY_DEBUG_FORMAT\e[0m" "$VERBOSITY_PREFIX" "DEBUG" "$1"
    echo -e "$message"
}