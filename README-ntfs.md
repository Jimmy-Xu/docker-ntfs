
# summary
- CREATE GPT
- create partition
  - ESP 100MB (mkfs fat32)
  - MSR 128MB
  - NTFS (mkfs ntfs)


# generate nanoserver DM device

```
//start dockerd
$ sudo ./dockerd -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock --storage-driver=devicemapper -g /mnt/sdc/var/lib/docker1.13-ntfs --storage-opt dm.fs=ntfs-3g --storage-opt dm.mkfsarg="-U" --storage-opt dm.mkfsarg="-p 2048" --storage-opt dm.mkfsarg="-f" --storage-opt dm.mountopt="offset=$((0x0e400000))"  --debug

//load nanoserver image
$ ./docker load -i ~/microsoft/nanoserver/nanoserver.tar.gz

//create nanoserver container
$ ./docker create --name nanoserver microsoft/nanoserver ping localhost -t

//start container
$ ./docker start nanoserver

//list DM device
$ ll /dev/mapper 
total 0
crw------- 1 root root 10, 236 Mar 10 08:49 control
lrwxrwxrwx 1 root root       7 Mar 11 15:52 docker-8:32-8388674-14a451b55ec3984d9c3230f269ed762a4866719d7e7862efca419ab8cd66e62a -> ../dm-1
lrwxrwxrwx 1 root root       7 Mar 11 15:53 docker-8:32-8388674-76e4b788a9fe892b3a3179d16cea08049d634bddd420eb6af6e84b287822d9ac -> ../dm-2
lrwxrwxrwx 1 root root       7 Mar 11 15:53 docker-8:32-8388674-f5da36d31e424b483a8f68a8db77f1aa129368784800522f62522042d34d3281 -> ../dm-4
lrwxrwxrwx 1 root root       7 Mar 11 15:53 docker-8:32-8388674-f5da36d31e424b483a8f68a8db77f1aa129368784800522f62522042d34d3281-init -> ../dm-3
lrwxrwxrwx 1 root root       7 Mar 11 15:51 docker-8:32-8388674-pool -> ../dm-0

```

# create NanoBoot.raw

```
//创建NanoBoot.raw
$ qemu-img create -f raw NanoBoot.raw 128M

//分区
$ sudo parted NanoBoot.raw --script \
mklabel gpt \
mkpart ESP fat32 0% 100% \
set 1 boot on \
print

//创建loop设备
$ LOOPDEV=$(losetup -f)
$ sudo losetup -o $((0x4400)) ${LOOPDEV} NanoBoot.raw

//查看loop设备
$ losetup ${LOOPDEV}
/dev/loop2: []: (/home/osboxes/iso/NanoBoot.raw), offset 17408


//格式化
$ sudo mkfs.fat -s1 -F32 ${LOOPDEV}

```

# qemu启动WinPE,并挂载NanoBoot.raw和DM device

```
IMG_PATH=${HOME}/iso
sudo qemu-system-x86_64 -enable-kvm -smp 1 -m 2048 \
  -bios /usr/share/edk2.git/ovmf-x64/OVMF_CODE-pure-efi.fd \
  -enable-kvm -netdev tap,id=t1,ifname=tap1,script=no,downscript=no -net nic,model=virtio,netdev=t1,macaddr=00:16:e4:9a:b3:6a \
  -machine vmport=off \
  -boot order=d,menu=on \
  -drive file=${IMG_PATH}/NanoBoot.raw,format=raw \
  -drive file=/dev/mapper/docker-8:32-8388674-f5da36d31e424b483a8f68a8db77f1aa129368784800522f62522042d34d3281,format=raw \
  -cdrom ${IMG_PATH}/Win8PE64.iso \
  -usb -usbdevice tablet \
  -vnc :8
  
调整：
1. DM device中
  1)复制rootfs/UtilityVM/Files下的所有目录到根目录
  2)复制rootfs/Files下所有文件到根目录(覆盖原有文件)
  3)复制rootfs/UtilityVM/Files/EFI到ESP分区
  4)WinPE中，用bootice修改ESP分区EFI/Microsoft/Boot/BCD
    a.修改ApplicationDevice和OSDevice参数为C:
    b.复制启动项
2. NanoBoot.raw中
  1)复制DM device的ESP中的EFI目录到当前ESP中
3. NanoGuestUEFI.qcow2中
  1)复制相关文件覆盖到DM device中
```

# 挂载NanoGuestUEFI.qcow2,NanoBoot.raw和DM device到文件系统

```
//挂载NanoGuestUEFI.qcow2(从Windows Server 2016提取的nanoserver)
sudo guestmount -a NanoGuestUEFI.qcow2 -m /dev/sda4 /mnt/ntfs

//挂载DM device(docker创建的nanoserver)
sudo mount -o loop,offset=$((0x100000)) /dev/mapper/docker-8:32-8388674-f5da36d31e424b483a8f68a8db77f1aa129368784800522f62522042d34d3281 /mnt/docker_esp
sudo mount -o loop,offset=$((0x0e400000)) /dev/mapper/docker-8:32-8388674-f5da36d31e424b483a8f68a8db77f1aa129368784800522f62522042d34d3281 /mnt/docker_ntfs

//挂载NanoBoot.raw(只有ESP)
sudo mount -o loop,offset=$((0x4400)) NanoBoot.raw /mnt/esp
```

