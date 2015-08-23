# jimmy-import
Merge multiple files in this format:

    # 22.08.2015
    * 06:15 5g creatine
    * 06:20 250mg caffeine

...into a single file in this format:

    supplements:
        - date: 2015-08-22T06:15:00Z
          name: Creatine
          dose: 5

        - date: 2015-08-22T06:20:00Z
          name: Caffeine
          dose: 0.25

## Usage
    jimmy-import example.md [example2.md] > example.yaml
