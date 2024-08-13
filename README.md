# kubetea

<img src="https://github.com/flyhope/kubetea/releases/download/v0.0.8/logo.jpeg" align="right" alt="LOGO" style="width: 200px; margin-right: 20px;"> 

<p> kubernetes simple cli ui client <p>
<p> 🎈 Easy to use</p>
<p> 🚀 Quick to get started</p>
<p> ⌨️ Use command-line UI instead of window UI</p>
<p> 💫 Improve work efficiency</p>

## Install

1. Download the latest release from the [releases page](https://github.com/flyhope/kubetea/releases) to you bin $PATH, and rename to `kubetea`.
2. cd your bin path and `chmod +x kubetea`.
3. fist start well write config to default config path `~/.kubetea/config.yaml`.

## Usage

### main ui
```bash
kubetea
```

key map (select line for table or page)

| Key        | Description           | 
|------------|-----------------------|
| j,↓        | table row select down |  
| k,↑        | table row select up   |      
| PageDown,f | page up               |
| PageUp,b   | page down             |
| d,ctrl+d   | page half down        |
| u,ctrl+u   | page half up          |
| G,end      | go to end             |
| g,home     | go to home            |
| 1-9        | sort by column index  |
| 0          | reset sort            |

key map (table)

| Key    | Description        | cluster  | pod | container |
|--------|--------------------|:--------:|:---:|:---------:|
| esc    | exit or go back    |    ✅️    | ✅️  |    ✅️     |
| ctrl+c | exit               |    ✅️    | ✅️  |    ✅️     |
| /      | search input focus |    ✅️    | ✅️  |    ✅️     |
| i      | pod infomation     |    ❌️    | ✅️  |    ✅️     |
| e      | pod describe       |    ❌️    | ✅️  |    ✅️     |
| s      | shell              |    ❌️    | ❌️  |    ✅️     |

page key map (exmaple: pod infomation)

| Key        | Description              |
|------------|--------------------------|
| /          | search input focus       |
| esc        | search input blur        |
| enter      | search mode              |
| n          | next match (search mode) |
| N          | prev match (search mode) |

### quick enter a pod container list

```bash
kubetea $ipOrPodname
# example
# kubetea 192.168.0.1
# kubetea nginx-pod-5b7f7d7d8-7j7j2
```

## Feature

- [x] support krs-auth

## Thanks

* kubectl https://github.com/kubernetes/kubectl
* bubbletea https://github.com/charmbracelet/bubbletea
* lingma https://tongyi.aliyun.com/lingma
