# fluorescences
a comic blog with private galleries

## Install Instructions

First initialize the internal secret

```
fluorescences --tenant user init secret
```

then create a user, this will generate a password

```
fluorescences --tenant user init user --username user
```

create the boilerplate data

```
fluorescences --tenant user init data
```

then start!

```
fluorescences --tenant user --address 0.0.0.0 --port 5000 start
```

or use the included systemd unit file
