


## Development

You must be installed go-i18n command

```bash
go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest
goi18n -help
```

More go-i18n document:
* english https://github.com/nicksnyder/go-i18n/blob/main/README.md
* chinese https://github.com/nicksnyder/go-i18n/blob/main/.github/README.zh-Hans.md

### Add a new language

example: add a new language `fr`

1. clone this repository and cd to it.
2. init your language file.
    ```bash
    export KUBETEA_LANG=fr
    goi18n extract -format=yaml -outdir="./lang"
    cd lang
    goi18n merge -format=yaml view.en.yaml  active.en.yaml
    touch translate.$KUBETEA_LANG.yaml
    goi18n merge -format=yaml active.en.yaml translate.$KUBETEA_LANG.yaml
    ```
3. translate your language file `translate.fr.yaml`.
4. update merge your language file `active.fr.yaml`.
    ```bash
   rename translate.$KUBETEA_LANG.yaml active.$KUBETEA_LANG.yaml
    ```

### Update existing language

1. cd to this work directory.
2. update different language translate file.
    ```bash
    goi18n extract -format=yaml -outdir="./lang"
    cd lang
    goi18n merge -format=yaml view.en.yaml  active.en.yaml
    goi18n merge -format=yaml active.*.yaml
    ```
3. Translate all the messages in the `translate.*.yaml` files.
4. Run command merge the translated messages into the active message files.
   ```bash
   goi18n merge -format=yaml active.*.yaml translate.*.yaml
   rm -f translate.*.yaml
   ```
   
### Add a new language file, but not build in kubetea

1. look up ↑ `Add a new language` and do it. 
2. copy `active.xx.yaml` to `~/.kubetea/lang/active.xx.yaml`.
