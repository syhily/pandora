# Blog Attachments Store

This is a repository for storing and uploading images to [Assetry](https://github.com/syhily/assetry).

## ⚙️ How to Install

* `make build`: Generate the executable file. Move the file into your `PATH`.
* `make install`: Install tool and set shell completion on macOS.

## 🧩 Usage

### 🪄 Tool Configuration

Sample config file

### 🖼️ Images Procession

Convert images, resize, rename and upload to Assetry with YAML metadata file.

```bash
pandora image -h
```

**Flags:**

```text
  -f, --format string   The image format (default "jpg")
      --height int      The optional image height, 0 for keep ratio
  -q, --quality int     The image quality
  -s, --source string   The image file path (absolute of relative)
  -t, --time string     The date time, in 20060102 format
  -u, --upload          Whether to upload image (default true)
  -w, --width int       The resized image width (default 1280)
```

### 🎵 Netease Music

Download and upload netease music to Assetry with YAML metadata file.

```bash
pandora music -h
```

**Flags:**

```text
  -i, --id int   The music ID
  -v, --vip      Download through 3rd netease VIP API
```

### ☁️ Upload Assets

Upload the file or directory to Assetry.

```bash
pandora upload -h
```

**Flags:**

```text
  -f, --file string     The file or directory to be upload
  -t, --target string   The target path for uploading
```

### ⏬ Backup Assets

Backup all the assets in Assetry to local directory.

```bash
pandora backup
```
