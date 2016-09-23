# Path-Tracer

A path tracer I wrote in order to learn golang.

Sample render:
2000 SPP, 5 Shadow rays per sample, 128 Adaptive samples. Rendered in 57 minutes.
![](http://i.imgur.com/FATwZgB.png)

Features:

- CPU based stochastic unidirectional path tracer
- Concurrent, uses all available cores
- Supports OBJ files
- Various material properties
- K-D tree acceleration
- Supports adaptive sampling 
- Thin lens model with depth of field effect
