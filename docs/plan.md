# **Plan**

###### **With Specific Files:**

`stevedore plan -f example.yaml -o example_override_example.yaml -e dev_environment.yaml`

###### **With directories:**
`stevedore plan -f services/ -o overrides/ -e envs/`

###### **What it does:**

Stevedore plan does a dry run over the manifests to be installed and shows the users the various
changes and new entities that are going to be installed/deleted. 

**Flags accepted:**

|  Flags 	| Description  	| Notes  	|
|---	|---	|---	|
| -f  	|  manifest file or directory with manifest files 	| Atleast one should be given |
| -o  	| override file or directory with override files  	| Optional |
| -e  	| env file or directory with env files  	| Optional |
| -a  	| output of artifacts path 	| Optional |
| -r  	| helm repo name to which the charts will be pushed/pulled  	| Optional. Defaults to `chartmuseum` |
| -t  	| helm timeout  	| Optional. Defaults to `600 seconds` |
| --recursive  	| respects `depends_on` in manifest files and recursively renders/plans/applies for those manifests too  	| Optional|
| --config  	|  config file that stevedore looks at for environment information 	| Optional. Defaults to `/usr/local/etc/stevedore-config.yaml` |
| --log-level  	|  sets the logger level 	| Optional. Defaults to `error` |


_**Atleast one manifest file should be given**_. Overrides/Envs file(s) are optional. If provided they will be used.

Stevedore goes through all the manifests, collates all matching overrides and envs for the appropriate manifests and show's the results of running the command to the user.
It shows the various kubernetes entities created/modified/deleted by the manifest(s) provided.

A very useful detail to be noted is the `depends_on` portion. Using this, one manifest can depend on another manifest.
If the recursive flag is provided, now stevedore will recursively also do a render/plan/apply for the manifests dependencies mentioned in this part of the manifest. 

This command ideally should precede the [apply](apply.md)  command. 

