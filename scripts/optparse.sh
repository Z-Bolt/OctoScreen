#!/usr/bin/env bash
declare -g optparse_license="$(cat <<EOL
# Optparse - a BASH wrapper for getopts
# https://github.com/nk412/optparse
# https://github.com/thebeline/optparse
#
# Copyright (c) 2015 Nagarjuna Kumarappan
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.
#
# Author:
#     Nagarjuna Kumarappan <nagarjuna.412@gmail.com>
#
# Contributors:
#     Alessio Biancone <alebian1996@gmail.com>
#     Mike Mulligan <mike@belineperspectives.com>
EOL
)"

function optparse.init(){
    
    unset optparse_version
    unset optparse_usage
    unset optparse_process
    unset optparse_defaults
    unset optparse_name
    unset optparse_usage_header
    unset optparse_variables_validate
    unset optparse_default_group
    
    declare -g optparse_version="0.0.2"
    
    declare -g optparse_defaults=""
    declare -g optparse_process=""
    declare -g optparse_variables_validate=""
    declare -g optparse_shortnames=()
    declare -g optparse_longnames=()
    declare -g optparse_variables=()
    
    for option in "$@"; do
        local key="${option%%=*}";
        local value="${option#*=}";

        case "$key" in
            "default_group"|"name"|"usage_header"|"description")
                declare -g optparse_$key="$value"
            ;;
            "help_full_width")
                if ! [[ "$value" =~ ^[0-9]+$ ]]; then
                        optparse.throw_error "optparse.init '$key' must be of type INTIGER"
                fi
                declare -g optparse_$key="$value"
            ;;
        esac
    done
    
    [ -z "$optparse_help_full_width" ] && declare -g optparse_help_full_width=80
    [ -z "$optparse_default_group" ] && declare -g optparse_default_group="OPTIONS"
    [ -z "$optparse_name" ] && declare -g optparse_name="$(basename $0)"
    [ -z "$optparse_usage_header" ] && declare -g optparse_usage_header="[OPTIONS]"
    [ -z "$optparse_description" ] && declare -g optparse_description="${optparse_name} Help"
    
    optparse.define long=help desc="This Help Screen" dispatch="optparse.usage" behaviour=flag help=explicit
    optparse.define long=optparse_license desc="The OptParse Library License" dispatch="optparse.license" help=hide
    
}

# -----------------------------------------------------------------------------------------------------------------------------
function optparse.throw_error(){
    local message="$1"
    [ ! -z $2 ] && message+=" for option: ($2)"
    optparse._log "ERROR" "$message"
    exit 1
}

function optparse._log(){
    local type="$1"
    local message="$2"
    echo "OPTPARSE $type: $message"
}

function optparse.warn(){
    local message="$1"
    [ ! -z $2 ] && message+=" for option: ($2)"
    optparse._log "WARN" "$message"
}

function optparse.group(){
    local group="$1"
    [ ! -z "$group" ] && [ ! -z "$optparse_usage" ] && group="#NL$group"
    optparse_usage+="echo \"$group\";#NL"
}

# -----------------------------------------------------------------------------------------------------------------------------
function optparse.define(){
    local errorname=""
    local short=""
    local shortname=""
    local long=""
    local longname=""
    local desc=""
    local default=""
    local behaviour="default"
    local list=false
    local variable=""
    local dispatch=""
    local val=""
    local has_val="false"
    local has_default="false"
    local behaviour_help="default"
    local behaviour_extra=""
    local optparse_help_indent=15
    local debug="false"

    for option in "$@"; do
        local key="${option%%=*}";
        local value="${option#*=}";
        
        case "$key" in
            "short")
                [ -z errorname ] &&
                    errorname="$value"
                [ ${#value} -ne 1 ] &&
                    optparse.throw_error "short name expected to be one character long" "$errorname"
                for i in "${optparse_shortnames[@]}"; do
                    if [ "$i" == "$value" ] ; then
                        optparse.warn "shortname [-$value] already handled" "$errorname"
                    fi
                done
                optparse_shortnames+=("$value")
                shortname="$value"
                short="-$value"
            ;;
            "long")
                [ -z ${value} ] &&
                    optparse.throw_error "long name expected to be atleast one character long" "$error_name"
                for i in "${optparse_longnames[@]}"; do
                    if [ "$i" == "$value" ] ; then
                        optparse.warn "longname [--$value] already handled" "$errorname"
                    fi
                done
                optparse_longnames+=("$value")
                longname="$value"
                errorname="$value"
                long="--$value"
            ;;
            "desc")
                desc="$value"
            ;;
            "default")
                default="$value"
                has_default="true"
            ;;
            "behaviour")
                case $value in
                    default|list|flag)
                        behaviour="$value"
                    ;;
                    *)
                        optparse.throw_error "behaviour [$value] not supported" "$errorname"
                    ;;
                esac
                ;;
            "list")
                list="$value"
            ;;
            "variable")
                variable="$value"
                for i in "${optparse_variables[@]}"; do
                    if [ "$i" == "$value" ] ; then
                        optparse.warn "value assignment [\$$value] already handled" "$errorname"
                    fi
                done
                optparse_variables+=("$value")
            ;;
            "value")
                val="$value"
            ;;
            "dispatch")
                dispatch="$value"
            ;;
            "explicit_help")
                always_help=false
            ;;
            "debug")
                debug=true
            ;;
            "extra"|"help")
                case $value in
                    default|explicit|hide)
                        declare behaviour_$key="$value"
                    ;;
                    *)
                        optparse.throw_error "$key [$value] not supported" "$errorname"
                    ;;
                esac
            ;;
        esac
    done
    
    [ -z $behaviour_extra ] && {
        [ "$behaviour" == "flag" ] && behaviour_extra="hide" || behaviour_extra="explicit"
    }
    
    [ -z "$errorname" ] && optparse.throw_error "argument must have a long or short name"
    
    flag=$([[ $behaviour == "flag" ]] && echo "true" || echo "false")
    is_list=$([[ $behaviour == "list" ]] && echo "true" || echo "false")

    [ $behaviour == "flag" ] && {
        [ -z $default ] && default=false
        has_default=true
        [ $default = "true" ] && val="false"
        [ $default = "false" ] && val="true"
    }
    
    has_val=$([[ -z "$val" ]] && echo "false" || echo "true")
    has_variable=$([[ -z "$variable" ]] && echo "false" || echo "true")

    # check list behaviour
    [ $behaviour == "list" ] && {
        [[ -z ${list:-} ]] &&
            optparse.throw_error "list is mandatory when using list behaviour" "$errorname"

        $has_default && {
            valid=false
            for i in $list; do
                [[ $default == $i ]] && valid=true && break
            done

            $valid || optparse.throw_error "default should be in list" "$errorname"
        }
    }

    if [ -z "$desc" ]; then
        if [ -z "$dispatch" ]; then
            optparse.throw_error "description is mandatory" "$errorname"
        else
            [ "$behaviour_help" == "default" ] && help_behaviour="explicit"
            [ "$behaviour_help" != "hide" ] && desc="Executes $dispatch"
        fi
    fi

    [ -z "$variable" ] && [ -z "$dispatch" ] && optparse.throw_error "you must give a target variable" "$errorname"
    
    
    [ -z "$optparse_usage" ] && optparse.group "$optparse_default_group"
    
    
    [ "$behaviour_help" != "hide" ] && {
        
        local _optparse_usage=""
    
        # build OPTIONS and help
        
        local description=()
        local description_index=0
        local description_word=""
        local description_sep=""
        local description_char_index=0
        
        # Break to lines
        
        #for (( description_char_index=0; description_char_index<${#desc}; description_char_index++ )); do
        #    #echo "${foo:$i:1}"
        #    :
        #done
        
        #_optparse_usage+="cat >&2 << EOU"
        _optparse_usage+="#TB$(printf "%-${optparse_help_indent}s %s" "${short:=  }$([ ! -z $short ] && [ ! -z $long ] && echo "," || echo " ") ${long}" "${desc}")"
        #        -h, --help      T
        
        _optparse_usage="cat >&2 << EOU#NL$_optparse_usage#NLEOU#NL"
        
        [ "$behaviour_extra" != "hide" ] && {
            
            local _optparse_extra=""
            
            $is_list && _optparse_extra+="#TB#TB#TBOne of: '$list'#NL"
        
            $flag && {
                _optparse_extra+="#TB#TB#TBTreated as a Flag#NL"
            } || {
                ${has_default} &&
                    _optparse_extra+="#TB#TB#TBDefault: '$default'#NL"
            }
            
            [ ! -z "$_optparse_extra" ] && {
            
                _optparse_extra="cat >&2 << EOE#NL${_optparse_extra}EOE#NL"
                
                [ "$behaviour_extra" == "explicit" ]  && 
                    _optparse_extra="if [[ \$1 == "true" ]]; then#NL$_optparse_extra#NLfi#NL"
                
                _optparse_usage+="$_optparse_extra"
            }
        }
        
        [ "$behaviour_help" == "explicit" ] && 
            _optparse_usage="if [[ \$1 == "true" ]]; then#NL$_optparse_usage#NLfi#NL"
    
        optparse_usage+="$_optparse_usage"
        
        #echo "$_optparse_usage"
    }
    
    $has_default && [ ! -z "$variable" ] &&
        optparse_defaults+="#NL${variable}='${default}'"
    
    local validate_variable="$is_list && {
                    valid=false
                    for i in $list; do
                        [[ \$optparse_processing_arg_value == \$i ]] && valid=true && break
                    done
                    \$valid || optparse.usage false \"ERROR: invalid value for argument \\\"\$_optparse_arg_errorname\\\"\" 1
                }
                $has_variable && [[ -z \${${variable}:-$($has_default && echo 'DEF' || echo '')} ]] && optparse.usage false \"ERROR: (\$_optparse_arg_errorname) requires input value\" 1"
        
    local dispatch_caller="# No Dispatcher"
    
    [ ! -z "$dispatch" ] && {
        [ -z "$dispatch_var" ] && {
            $has_variable && dispatch_var="\"\$${variable}\"" || dispatch_var="\"\$optparse_processing_arg_value\""
        }
        dispatch_caller="${dispatch} ${dispatch_var}"
    }
    
    [ ! -z "$shortname" ] && {
        optparse_process_short+="
            ${shortname})
                _optparse_arg_errorname=\"$short\"
                optparse_processing_arg_value=\"\"
                ( $flag || $has_val ) && optparse_processing_arg_value=\"$val\"  || {
                    $has_variable && {
                        optparse_processing_arg_value=\"\$optparse_processing_arg\";
                        if [ -z \"\$optparse_processing_arg_value\" ]; then
                            optparse_processing_arg_value=\"\$1\" && shift
                        else
                            optparse_processing_arg=''
                        fi
                    }
                }
                $has_variable && ${variable}=\"\$optparse_processing_arg_value\"
                $validate_variable
                $dispatch_caller
                continue
            ;;"
    }
    
    [ ! -z "$longname" ] && {
        optparse_process_long+="
            ${longname})
                _optparse_arg_errorname=\"$long\"
                $has_val && {
                    [ ! -z \"\$optparse_processing_arg_value\" ] && optparse.usage true 'ERROR: (${errorname}) does not accept user input' 1
                    optparse_processing_arg_value=\"$val\"
                }
                $has_variable && ${variable}=\"\$optparse_processing_arg_value\"
                $validate_variable
                $dispatch_caller
                continue
            ;;"
    }
    
    $has_variable && optparse_variables_validate+="
        [[ -z \${${variable}:-$($has_default && echo 'DEF' || echo '')} ]] && optparse.usage true 'ERROR: (${errorname}) not set' 1 || true"
}

# -----------------------------------------------------------------------------------------------------------------------------
function optparse.build(){
    
    local preserve_positional=false
    local allow_unrecognized=false
    local optparse_debug=false
    local unknown_long_handler=''
    local unknown_short_handler=''
    local non_empty_shorts=''
    
    local _optparse_process=$optparse_process
    local _optparse_variables_validate=$optparse_variables_validate

    for option in "$@"; do
        local key="${option%%=*}";
        local value="${option#*=}";

        case "$key" in
            "version")
                optparse_version="$value"
            ;;
            "allow")
                value=( $value )
                for allowed in "${value[@]}"; do
                    case "$allowed" in
                        "positional")
                            preserve_positional=true
                        ;;
                        "unrecognized")
                            allow_unrecognized=true
                        ;;
                        *)
                            optparse.throw_error "Unhandled Allowance '$allow'"
                        ;;
                    esac
                done
            ;;
            "debug")
                optparse_debug=true
            ;;
        esac
    done
    
    $preserve_positional && {
        preserve_positional="
        [ ! -z \"\$optparse_processing_arg\" ] && {
            $optparse_debug && echo \"Passing positional arg to args \\\"\$optparse_processing_arg\\\"\"
            optparse_processing_args+=( \"\$optparse_processing_arg\" )
        }"
    } || {
        preserve_positional="optparse.usage true \"ERROR: Unconfigured Arguments are not accepted.\" 1"
    }
    
    $allow_unrecognized && {
        unknown_long_handler="
                $optparse_debug && echo \"Passing unknown long to args \$optparse_processing_arg\"
                optparse_processing_args+=( \"\$optparse_processing_arg\" )
                continue"
        unknown_short_handler="
                $optparse_debug && echo \"Passing unknown short '\$optparse_processing_arg_key' to args \$optparse_processing_arg\"
                optparse_processing_arg+=\"\$optparse_processing_arg_key\"
                continue"
        non_empty_shorts="
            [ ! -z \"\$optparse_processing_arg\" ] && {
                $optparse_debug && echo \"Passing non-empty shorts '-\$optparse_processing_arg' to args \$optparse_processing_arg\"
                optparse_processing_args+=(\"-\$optparse_processing_arg\")
                continue
            }"
    } || {
        unknown_long_handler="optparse.usage true \"Unrecognized option: \$optparse_processing_arg_key\" 1"
        unknown_short_handler="$unknown_long_handler"
    }
    
    # Function usage
    cat <<EOF | sed -e 's/#NL/\n/g' -e 's/#TB/\t/g'
function optparse.license(){ cat <<EOL
$optparse_license
EOL
exit 0
}

function optparse.usage(){
( [ "\$1" == "true" ] && [ ! -z "\$2" ]) && echo -e "\$2"
cat >&2 << EOH
${optparse_description}
Usage: $optparse_name $optparse_usage_header

EOH
$optparse_usage
cat >&2 << EOH

\$(printf "%${optparse_help_full_width}s %s" "Powered by optparse $optparse_version")
\$(printf "%${optparse_help_full_width}s %s" "@see --optparse_license")
EOH
if [ "\$1" == "true" ] && [ ! -z "\$3" ]; then
    exit "\$3"
else
    exit 3
fi
}

# Set default variable values
$optparse_defaults

optparse_processing_args=()

# Begin Parseing
while [ \$# -ne 0 ]; do
    optparse_processing_arg="\$1"
    shift
    
    case "\$optparse_processing_arg" in
        --)
            optparse_processing_args+=("--") && break
            ;;
            
        --*)
            
            optparse_processing_arg_key="\${optparse_processing_arg%%=*}";
            optparse_processing_arg_value="\${optparse_processing_arg#*=}";
            
            [ "\$optparse_processing_arg_value" == "\$optparse_processing_arg_key" ] && optparse_processing_arg_value=''
            
            case "\${optparse_processing_arg_key:2}" in
                $optparse_process_long
            esac
            
            $unknown_long_handler
            
            ;;
            
        -*)
            
            optparse_processing_arg_short="\${optparse_processing_arg:1}"
            optparse_processing_arg=''
            
            while [ \${#optparse_processing_arg_short} -ne 0 ]; do
            
                optparse_processing_arg_key="\${optparse_processing_arg_short:0:1}";
                optparse_processing_arg_short="\${optparse_processing_arg_short:1}"
                optparse_processing_arg_value='';
                
                case "\$optparse_processing_arg_key" in
                    $optparse_process_short
                esac
                
                $unknown_short_handler
                
            done
            
            $non_empty_shorts
            
            ;;
        
        $default_handler
        
    esac
    
    $preserve_positional
    
done

# GODLY arg quote with printf, WHERE HAVE YOU BEEN??
eval set -- "\$([ ! -z "$optparse_processing_args" ] && printf ' %q' "\${optparse_processing_args[@]}") \$(printf ' %q' "\${@}")"

$optparse_variables_validate

unset _optparse_params
unset _optparse_param

EOF

    if [ -z "$optparse_preserve" ] ; then
        # Unset global variables
        optparse.init
    fi
}
optparse.init