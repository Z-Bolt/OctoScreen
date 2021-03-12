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
    
    declare -g optparse_version="0.0.2"
    
    declare -g optparse_defaults=""
    declare -g optparse_process=""
    declare -g optparse_variables_validate=""
    declare -g optparse_shortnames=()
    declare -g optparse_longnames=()
    declare -g optparse_variables=()
    declare -g optparse_usage_commands=""
    
    for option in "$@"; do
        local key="${option%%=*}";
        local value="${option#*=}";

        case "$key" in
            "default_group"|"usage_header"|"description")
                declare -g optparse_$key="$value"
            ;;
            "name")
                [[ "${value::1}" == "-" ]] &&
                    declare -g optparse_name="" ||
                    declare -g optparse_name="$(basename $value)"
            ;;
            "help_full_width")
                if ! [[ "$value" =~ ^[0-9]+$ ]]; then
                        optparse.throw_error "optparse.init '$key' must be of type INTIGER"
                fi
                declare -g optparse_$key="$value"
            ;;
        esac
    done
    
    [[ -z "$optparse_name" ]] && {
        [[ "${0::1}" == "-" ]] &&
            declare -g optparse_name="" ||
            declare -g optparse_name="$(basename $0)"
    }
    [[ -z "$optparse_help_full_width" ]] && declare -g optparse_help_full_width=80
    [[ -z "$optparse_default_group" ]] && declare -g optparse_default_group="OPTIONS"
    [[ -z "$optparse_usage_header" ]] && declare -g optparse_usage_header="[OPTIONS]"
    [[ -z "$optparse_description" ]] && declare -g optparse_description="${optparse_name} Help"
    
    optparse.define long=help desc="This Help Screen" dispatch="optparse.usage" behaviour=flag help=explicit
    optparse.define long=optparse_license desc="The OptParse Library License" dispatch="optparse.license" help=hide
    
}

function optparse.reset(){
    optparse.unset
    optparse.init $(printf ' %q' "${@}")
}

function optparse.unset(){
    unset optparse_version
    unset optparse_usage
    unset optparse_process
    unset optparse_defaults
    unset optparse_name
    unset optparse_usage_header
    unset optparse_variables_validate
    unset optparse_default_group
}

# -----------------------------------------------------------------------------------------------------------------------------
function optparse.throw_error(){
    local message="$1"
    [[ ! -z $2 ]] && message+=": ($2)"
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
    [[ ! -z "$2" ]] && message+=": ($2)"
    optparse._log "WARN" "$message"
}

function optparse.group(){
    local group="$1"
    [[ ! -z "$group" && ! -z "$optparse_usage" ]] && group="#NL$group"
    optparse_usage+="cat >&2 << EOU#NL$group#NLEOU#NL"
}

# -----------------------------------------------------------------------------------------------------------------------------
function optparse.define.single(){
    local name=""
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
        
        [[ $value == $key ]] && value=''
        
        case "$key" in
            desc|description)
                desc="$value"
            ;;
            name)
                name="$value"
            ;;
            value)
                val="$value"
            ;;
            dispatch)
                dispatch="$value"
            ;;
            help)
                case $value in
                    default|explicit|hide)
                        declare behaviour_$key="$value"
                    ;;
                    *)
                        optparse.throw_error "$key [$value] not supported" $name
                    ;;
                esac
            ;;
        esac
    done
    
    [[ -z "$val" ]] &&
        local has_val=false || 
        local has_val=true;
    
    [[ -z "$name" ]] &&
        optparse.throw_error "name is mandatory";
        
    [[ -z "$desc" && "$behaviour_help" != "hide" ]] &&
        optparse.throw_error "description is mandatory" "$name";
        
    [[ -z "$dispatch" ]] &&
        optparse.throw_error "a dispatcher is mandatory for commands" "$name";
    
    #[ -z "$optparse_usage" ] && optparse.group "$optparse_default_group"
    
    
    [[ $behaviour_help != hide ]] && {
        
        [[ -z "$optparse_usage_commands" ]] &&
            optparse_usage_commands="cat >&2 << EOU#NLCommands#NLEOU#NL"
        
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
        _optparse_usage+="#TB$(printf "%-${optparse_help_indent}s %s" "   ${name}" "${desc}")"
        #        -h, --help      T
        
        _optparse_usage="cat >&2 << EOU#NL$_optparse_usage#NLEOU#NL"
        
        [ "$behaviour_help" == "explicit" ] && 
            _optparse_usage="if [[ \$1 == "true" ]]; then#NL$_optparse_usage#NLfi#NL"
    
        optparse_usage_commands+="$_optparse_usage"
        
        #echo "$_optparse_usage"
    }
    
    optparse_additional_handlers+="
        ${name})
            ${dispatch} ; true;
            ;;"
    
    true
}

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
    local optional="false"
    local has_variable="false"

    for option in "$@"; do
        local key="${option%%=*}";
        local value="${option#*=}";
        
        case "$key" in
            short)
                [[ ${#value} -ne 1 ]] &&
                    optparse.throw_error "short name expected to be one character long" $value;
                
                for i in "${optparse_shortnames[@]}"; do
                    [[ $i == $value ]] && 
                        optparse.warn "shortname [-$value] already handled" "-$value" && break;
                done;
                
                local optparse_shortnames+=("$value")
                local shortname="$value"
                local short="-$value"
                [[ -z "$errorname" ]] &&
                    local errorname="$short";
            ;;
            long)
                [[ -z ${value} ]] &&
                    optparse.throw_error "long name expected to be atleast one character long" "--$value";
                
                for i in "${optparse_longnames[@]}"; do
                    [[ $i == $value ]] &&
                        optparse.warn "longname [--$value] already handled" "--$value" && break;
                done;
                
                optparse_longnames+=("$value")
                longname="$value"
                long="--$value"
                errorname="$long"
            ;;
            desc|description)
                desc="$value"
            ;;
            default)
                default="$value"
                has_default=true
            ;;
            optional)
                optional=true
            ;;
            flag|list)
                behaviour=$key
            ;;
            behaviour)
                case $value in
                    default|list|flag)
                        behaviour=$value
                    ;;
                    *)
                        optparse.throw_error "behaviour [$value] not supported" $errorname
                    ;;
                esac;
                
                ;;
            list)
                list="$value"
            ;;
            variable)
                local variable=$value
                
                for i in "${optparse_variables[@]}"; do
                    [[ $i == $value ]] &&
                        optparse.warn "value assignment [\$$value] already handled" $errorname && break;
                done;
                
                optparse_variables+=("$value")
            ;;
            value)
                val=$value
            ;;
            dispatch)
                dispatch=$value
            ;;
            extra|help|hide|explicit)
                case $value in
                    default|explicit|hide)
                        [[ $value == $key ]] && {
                            local behaviour_help=$value;
                            local behaviour_extra=$value;
                        } || 
                            local behaviour_$key=$value;
                    ;;
                    *)
                        optparse.throw_error "$key [$value] not supported" "$errorname"
                    ;;
                esac
            ;;
        esac
    done
    
    [[ -z $behaviour_extra ]] && {
        [[ $behaviour == flag ]] &&
            local behaviour_extra=hide ||
            local behaviour_extra=explicit;
    }
    
    [[ -z "$errorname" ]] && optparse.throw_error "argument must have a long or short name"
    
    [[ $behaviour == flag ]] &&
        local flag=true ||
        local flag=false;
    
    [[ $behaviour == list ]] &&
        local is_list=true ||
        local is_list=false;

    [[ $behaviour == flag ]] && {
        [[ -z $default ]] && default=false
        has_default=true
        [[ $default == true ]] && val=false
        [[ $default == false ]] && val=true
    }
    
    [[ -z "$val" ]] &&
        local has_val="false" ||
        local has_val="true";
    
    
    [[ -z "$variable" ]] &&
        local has_variable=false ||
        local has_variable=true;

    # check list behaviour
    [[ $behaviour == list ]] && {
        [[ -z ${list:-} ]] &&
            optparse.throw_error "list is mandatory when using list behaviour" $errorname

        $has_default && {
            local valid=false;
            for i in $list; do
                [[ $default == $i ]] &&
                    local valid=true &&
                    break;
            done;

            $valid || optparse.throw_error "default should be in list" $errorname
        }
    }

    [[ -z "$desc" && $behaviour_help != hide ]] && {
        [[ -z "$dispatch" ]] && {
            optparse.throw_error "description is mandatory" $errorname
        } || {
            [[ $behaviour_help == default ]] && help_behaviour=explicit
            [[ $behaviour_help != hide ]] && desc="Executes $dispatch"
        }
    }

    [[ -z "$variable" && -z "$dispatch" ]] && optparse.throw_error "you must give a target variable" "$errorname";
    
    
    [[ $behaviour_help != hide ]] && {
        
        [[ -z "$optparse_usage" ]] && optparse.group "$optparse_default_group";
        
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
        
        [[ $behaviour_extra != hide ]] && {
            
            local _optparse_extra=""
            
            $is_list && _optparse_extra+="#TB#TB#TBOne of: '$list'#NL"
        
            $flag && {
                _optparse_extra+="#TB#TB#TBTreated as a Flag#NL"
            } || {
                $has_default &&
                    _optparse_extra+="#TB#TB#TBDefault: '$default'#NL"
            }
            
            [[ ! -z "$_optparse_extra" ]] && {
            
                _optparse_extra="cat >&2 << EOE#NL${_optparse_extra}EOE#NL"
                
                [[ $behaviour_extra == explicit ]]  && 
                    _optparse_extra="if [[ \$1 == "true" ]]; then#NL$_optparse_extra#NLfi#NL"
                
                _optparse_usage+=$_optparse_extra
            }
        }
        
        [[ $behaviour_help == explicit ]] && 
            _optparse_usage="if [[ \$1 == true ]]; then#NL$_optparse_usage#NLfi#NL"
    
        optparse_usage+=$_optparse_usage
        
        #echo "$_optparse_usage"
    }
    
    $has_default && [[ ! -z "$variable" ]] && 
        optparse_defaults+="#NL[[ -z \$${variable} ]] && ${variable}='${default}'"
    
    local validate_variable="$is_list && {
                    valid=false
                    for i in $list; do
                        [[ \$__arg_value == \$i ]] && valid=true && break
                    done
                    \$valid || optparse.usage false \"ERROR: invalid value for argument \\\"\$__arg_errorname\\\"\" 1
                }
                $has_variable && [[ -z \${${variable}:-$($has_default && echo 'DEF' || echo '')} ]] && optparse.usage false \"ERROR: (\$__arg_errorname) requires input value\" 1 || true"
        
    local dispatch_caller="# No Dispatcher"
    
    [[ ! -z "$dispatch" ]] && {
        [[ -z "$dispatch_var" ]] && {
            $has_variable &&
                dispatch_var="\"\$${variable}\"" ||
                dispatch_var="\"\$__arg_value\""
        }
        dispatch_caller="${dispatch} ${dispatch_var}; true;"
    }
    
    [[ ! -z "$shortname" ]] && {
        optparse_process_short+="
            ${shortname})
                __arg_errorname=$short;
                __arg_value='';
                ( $flag || $has_val ) && {
                    __arg_value=\"$val\";
                }
                $has_variable && [[ -z \"\$__arg_value\" ]] && {
                    __arg_value=$__arg_sremain;
                    __arg_sremain='';
                    [[ -z \"\$__arg_value\" && ! -z \"\$1\" ]] && __arg_value=\"\$1\" && shift;
                }
                $has_variable && ${variable}=\"\$__arg_value\";
                $validate_variable
                $dispatch_caller
                continue
            ;;"
    }
    
    [[ ! -z "$longname" ]] && {
        optparse_process_long+="
            ${longname})
                __arg_errorname=$long;
                $has_val && {
                    [[ ! -z \"\$__arg_value\" ]] && optparse.usage true 'ERROR: (${errorname}) does not accept user input' 1;
                    __arg_value=\"$val\";
                }
                $has_variable && ${variable}=\"\$__arg_value\";
                $validate_variable
                $dispatch_caller
                continue
            ;;"
    }
    
    $has_default &&
        local _def=DEF ||
        local _def='';
    
    $has_variable && optparse_variables_validate+="
        $optional || { [[ -z \${${variable}:-$_def} ]] && optparse.usage true 'ERROR: (${errorname}) not set' 1 || true; };"
    
    true
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
    
    local _exec=''

    for option in "$@"; do
        local key="${option%%=*}";
        local value="${option#*=}";
        
        case "$key" in
            name|description|usage_header|usage)
                declare -g optparse_$key="$value"
            ;;
            allow)
                value=( $value )
                for allowed in "${value[@]}"; do
                    case "$allowed" in
                        positional)
                            local preserve_positional=true
                        ;;
                        unrecognized)
                            local allow_unrecognized=true
                        ;;
                    esac
                done
            ;;
            finally)
                _exec="#NL$value;"
            ;;
        esac
    done
    
    [[ -z "$optparse_additional_handlers" ]] &&
        local _has_additional=false || 
        local _has_additional=true;
    
    $_has_additional &&
        local preserve_positional=true
    
    $preserve_positional && {
        preserve_positional="
        [ ! -z \"\$__arg\" ] && __args_processed+=( \"\$__arg\" )"
    } || {
        preserve_positional="optparse.usage true \"ERROR: Unconfigured Arguments are not accepted.\" 1"
    }
    
    $allow_unrecognized && {
        unknown_long_handler="
                __args_processed+=( \"\$__arg\" )
                continue"
        unknown_short_handler="
                __arg_processed+=\"\$__arg_key\"
                continue"
        non_empty_shorts="
            [ ! -z \"\$__arg_short_processed\" ] && {
                __args_processed+=(\"-\$__arg_processed\")
                continue
            }"
    } || {
        unknown_long_handler="optparse.usage true \"Unrecognized option: \$__arg_key\" 1"
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
[[ ! -z "\$2" ]] && echo -e "\$2"
cat >&2 << EOH
${optparse_description}
Usage: $optparse_name $optparse_usage_header

EOH
$optparse_usage_commands
$optparse_usage
cat >&2 << EOH

\$(printf "%${optparse_help_full_width}s %s" "Powered by optparse $optparse_version")
\$(printf "%${optparse_help_full_width}s %s" "@see --optparse_license")
EOH


[[ "\$1" == "true" && ! -z "\$3" ]] && exit "\$3" || exit 3
}



# Set default variable values
$optparse_defaults

__args_processed=()

optparse_command_hit="false"

# Begin Parseing
while [[ \$# -ne 0 ]]; do
    __arg_key=''
    __arg_value=''
    __arg_short=''
    __arg="\$1"
    shift
    
    case "\$__arg" in
        --)
            __args_processed+=("--") && break
            ;;
            
        --*)
            
            __arg_key="\${__arg%%=*}";
            __arg_value="\${__arg#*=}";
            
            [[ "\$__arg_value" == "\$__arg_key" ]] && __arg_value=''
            
            case "\${__arg_key:2}" in
                $optparse_process_long
            esac
            
            $unknown_long_handler
            
            ;;
            
        -*)
            
            __arg_sremain="\${__arg:1}"
            __arg_processed=''
            
            while [[ \${#__arg_sremain} -ne 0 ]]; do
            
                __arg_key="\${__arg_sremain:0:1}";
                __arg_sremain="\${__arg_sremain:1}"
                __arg_value='';
                
                case "\$__arg_key" in
                    $optparse_process_short
                esac
                
                $unknown_short_handler
                
            done
            
            $non_empty_shorts
            
            ;;
            
        $optparse_additional_handlers
        
    esac
    
    $preserve_positional
    
done

# GODLY arg quote with printf, WHERE HAVE YOU BEEN MY WHOLE LIFE??
[[ ! -z "\$__args_processed" ]] && {
    __args_processed="\$(printf '%q ' "\${__args_processed[@]}")"
    [[ ! -z "\$@" ]] && __args_processed+="\$(printf '%q ' "\${@}")"
    eval set -- "\$__args_processed"
}

$optparse_variables_validate

unset _optparse_params
unset _optparse_param

$_exec

EOF

    # Unset global variables
    [[ -z "$optparse_preserve" ]] && optparse.unset
}

optparse.init && true
