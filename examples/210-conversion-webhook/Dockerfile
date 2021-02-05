FROM flant/shell-operator:latest
LABEL app="shell-operator" example="210-conversion-webhook"

ADD framework/frameworks/shell/hook.sh /frameworks/shell/hook.sh

# Add hooks
ADD hooks /hooks

# Compiled version of shell-operator binary
ADD shell-operator /
