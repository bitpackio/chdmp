# chdmp
chunk based stream dump

## Usage
chdmp -input="pbx01_vda" -output="pbx01_vda.bak" -force -hash -verbose

8589934592 Bytes readed, 51232768 Bytes writen, 1048576 Chunks readed, 6254 Chunks writen, 1m17.12200587s, 111.38 MB/s

# vmsnap
save kvm domain

## Usage
vmsnap -snapsize=4 -snapstore="/srv/lib/libvirt/backups" -domain=pbx01 -disks=vda -dump -state -force -hash -verbose 

domain pbx01: detected target vda

domain pbx01: detected source /dev/vg00/pbx01_root

domain pbx01: using snapshot /dev/vg00/pbx01_root_2bf1a995-e585-46af-854c-d4076e91e7c3

domain pbx01: using destination /srv/lib/libvirt/backups/pbx01_vda

domain pbx01: saving state to /srv/lib/libvirt/backups/pbx01.state

domain pbx01: state saved to /srv/lib/libvirt/backups/pbx01.state

domain pbx01: creating snapshot /dev/vg00/pbx01_root_2bf1a995-e585-46af-854c-d4076e91e7c3

domain pbx01: restoring state from /srv/lib/libvirt/backups/pbx01.state

domain pbx01: state restored from /srv/lib/libvirt/backups/pbx01.state

domain pbx01: updating chunks from vda to /srv/lib/libvirt/backups/pbx01_vda

domain pbx01: 8589934592 Bytes readed, 2315370 Bytes writen, 2096129 Chunks readed, 565 Chunks writen, 1m40.345312376s, 85.6 MB/s 

domain pbx01: chunks from vda to /srv/lib/libvirt/backups/pbx01_vda updated

domain pbx01: removing snapshot /dev/vg00/pbx01_root_2bf1a995-e585-46af-854c-d4076e91e7c3

# Setup
Install golang packages.

apt-get install golang

go build src/chdmp/chdmp.go

go build src/vmsnap/vmsnap.go
