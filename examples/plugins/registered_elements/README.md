# Registered Elements example

This example shows how you can define custom gstreamer elements in go, register them in the element factory and use them in the same application in a pipeline.

We define two elements:

* `gocustombin` a custom GstBin that uses an audiomixer to aggregate the input of two `gocustomsrc`
* `gocustomsrc` a custom GstBin that uses an audiotestsrc and a volume element.