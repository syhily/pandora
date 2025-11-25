# Blog Attachments Store

This is a repository for storing and uploading images to [S3](https://cdn.yufan.me).

## ⚙️ How to Install

* `make build`: Generate the executable file. Move the file into your `PATH`.
* `make install`: Install tool and set shell completion on macOS.

## 🧩 Usage

### 🪄 Tool Configuration

Interactive configure pandora tool.

```bash
pandora config -h
```

### 🖼️ Convert Images

Convert images, resize, rename and upload to S3.

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

### ✂️ Split Large Chinese Font

Split large TTF fonts and deploying them into S3.

```bash
pandora font -h
```

**Flags:**

```text
  -t, --ttf string   The font path
```

### 🎵 Download Netease Music

Download and upload Netease music to S3 with metadata file.

```bash
pandora music -h
```

**Flags:**

```text
  -i, --id int   The music ID
  -v, --vip      Download through VIP API
```

### ☁️ Upload Attachments

Sync files to S3. Generate image metadata file with thumbhash.

```bash
pandora sync
```
