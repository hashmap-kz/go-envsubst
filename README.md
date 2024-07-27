# Architecture

- [Project postulated](#project-postulates)
- [Configuration](#configuration)

## Project postulates:

### Simplicity

> **Everywhere: in design, in configuration, in managing**
>
> The reason for that: is you need to be sure what is going on when you start your application.
>
> You have to be sure what this application does, and how to handle all errors as quickly as possible.
>
> It's not a goal to do 1000 things, but here is a main goal: make a few things right.
>
> This project is oriented specifically on Unix based OS, and it will be.

## Configuration

```
# you may use comma or spaces, but not both at the same time

# expand ONLY these list of env-vars and 
# ignore anything else (this is important to not expand nginx conf, for example)
#
export GENVSUBST_ALLOWED='APP_NAME, IMAGE_NAME, IMAGE_TAG'

# do NOT expand these env-vars, never. 
# this setting has the highest priority
#
export GENVSUBST_RESTRICTED='SECRET_X SECRET_Y'

# expand ONLY if the name of env-var starts with any of these prefixes
#
export GENVSUBST_ALLOWED_WITH_PREFIX='CI_ APP_'

# do NOT expand any of env-var if its name starts with any of these prefixes. 
# this setting is also has the highest priority.
#
export GENVSUBST_RESTRICTED_WITH_PREFIX='SECRET_ VAULT_'
```

> The combination is as simple as it looks like:
>
> Let's assume, you set all the possible options at the same time.
>
> This pseudocode shows the logic:

```
if (variable in GENVSUBST_ALLOWED) || 
    starts_with(variable, any(GENVSUBST_ALLOWED_WITH_PREFIX)) 
{
    
    // this means that GENVSUBST_ALLOWED or GENVSUBST_ALLOWED_WITH_PREFIX is set, and 
    // we have to check whether the variable is allowed for being expanded
    //
    allowToExpand = (variable not in GENVSUBST_RESTRICTED) 
        && !starts_with(variable, any(GENVSUBST_RESTRICTED_WITH_PREFIX))
    if allowToExpand {
        expand(variable)
    }
} else {

    // this means that GENVSUBST_ALLOWED or GENVSUBST_ALLOWED_WITH_PREFIX is not set, and 
    // we have to check whether the variable is allowed for being expanded
    //
    allowToExpand = (variable not in GENVSUBST_RESTRICTED) 
        && !starts_with(variable, any(GENVSUBST_RESTRICTED_WITH_PREFIX))
    if allowToExpand {
        expand(variable)
    }    
}
```







































