%%%
title = "Pod Conversion Conventions"
abbrev = "Pod"
category = "info"
area = "Koki Short"
keyword = ["Pod"]

[[author]]
fullname = "Caascade Labs, Inc."
%%%

.# Abstract

This document provides the details of the conversion mechanisms between Kubernetes Pod syntax and Koki Short syntax

{mainmatter}

# Introduction {#introduction}

A Pod is the unit of execution when using the Kubernetes orchestrator. It is defined in the api group core/v1. This version of Koki Short supports all valid Pod definitions of Kubernetes Versions \[v1.0 - v1.8]

A> No data is lost when converting between koki and kubernetes types

# Conversion from Kubernetes to Koki 

The Kubernetes API Object for Pod has a kind Key that always refers to value "Pod". This information is redundant, since Kubernetes parsers identify the type, and provide us with the Pod object, therefore this Kind field is removed in Koki and instead a new top level key called pod is used to denote a pod.

A kubernetes pod definitions looks like

~~~
   apiVersion: v1
   kind: pod
   metadata:
     name: pod_name
   ...
~~~

A koki pod definition written using Short syntax would look like 

~~~
    pod: 
      name: pod_name
    ...
~~~

As you can see, without getting into the most reductive parts, pod definitions in koki short syntax already look cleaner.

# Conversion from Koki to Kubernetes



# Examples

{backmatter}
