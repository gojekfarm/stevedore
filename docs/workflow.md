---
id: workflow
title: Workflow
sidebar_label: Workflow
---
###### **How a release is created with stevedore:**

Stevedore manages releases with 4 commands. These commands can be executed separately in any order at any time. But the preferred order in which they are to be run is -   

- [Init](init.md) - if no tiller is installed in the namespace or if the tiller version installed is not compatible with requirements of stevedore.
- [Render](render.md) - to render the various manifests with collated values to be replaced in the manifests.
- [Plan](plan.md) - to show the various kubernetes entities that will eb created/modified/deleted.
- [Apply](apply.md) - to perform the operations shown in the plan stage.

 

Todo: Insert flowchart/image here for the flow
