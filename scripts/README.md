# Document classifier (Python)

Used by sync_server to detect document-like images (whiteboard, notebook, book pages) and move them to Trash when `SYNC_DOCUMENT_TO_TRASH=1` and `SYNC_DOCUMENT_CLASSIFIER_PATH` point here.

## Usage

- **Heuristic only** (no ML): needs Python 3 and Pillow.
  ```bash
  pip install Pillow
  ```
  Set env: `SYNC_DOCUMENT_CLASSIFIER_PATH=scripts/document_classifier.py` (relative to sync_server exe dir) or full path.

- **With ONNX model**: install `onnxruntime` and set `DOCUMENT_CLASSIFIER_MODEL` to your `.onnx` file path. Model input: image (e.g. 224×224 RGB), output: 2 classes (document, photo) or single score. If the model file is missing or inference fails, the script falls back to the heuristic.

## Run manually

```bash
python scripts/document_classifier.py /path/to/image.jpg
```
Prints `document` or `photo` to stdout.