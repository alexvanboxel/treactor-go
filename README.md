# Reactor

##

PORT=3330
NAME=reactor
VERSION=1
MODE=local
GOOGLE_APPLICATION_CREDENTIALS=/usr/local/etc/keys/my-project.json
GOOGLE_PROJECT_ID=my-project


### Molecule spec


```
S    [H,x=1,y=2,z=3]*2[O,x=1,y=2,z=3],x=1,y=2,z=3
A     H,x=1,y=2,z=3
A                      O,x=1,y=2,z=3
```


```
S    2[5[Ur,log:1,xyz:4]^5[C,log:1,xyz:4]],x:1,y:2
O1     5[Ur,log:1,xyz:4]^5[C,log:1,xyz:4]
A        Ur,log:1,xyz:4
A                          C,log:1,xyz:4
```
