# Copyright 2026 Take Control - Software & Infrastructure
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Document vs photo classifier. Called by sync_server with one argument: image path.
# Prints "document" or "photo" to stdout. Uses heuristic by default; can load an
# ONNX model from DOCUMENT_CLASSIFIER_MODEL env var for better accuracy.

import os
import sys
from pathlib import Path

# Heuristic: same logic as Go (brightness + bimodal)
def looks_like_document_heuristic(image_path: str) -> bool:
    try:
        from PIL import Image
    except ImportError:
        return False
    try:
        img = Image.open(image_path).convert("RGB")
    except Exception:
        return False
    w, h = img.size
    if w < 10 or h < 10:
        return False
    # Resize to 200px width
    img.thumbnail((200, 200 * h // max(w, 1)))
    w, h = img.size
    pixels = w * h
    total = 0
    light = 0
    dark = 0
    for (r, g, b) in img.getdata():
        bright = (r + g + b) // 3
        total += bright
        if bright >= 240:
            light += 1
        elif bright <= 25:
            dark += 1
    mean = total // pixels
    ratio = (light + dark) / pixels
    return mean >= 140 and ratio >= 0.35


def main():
    if len(sys.argv) < 2:
        print("photo")
        sys.exit(0)
    image_path = sys.argv[1]
    if not os.path.isfile(image_path):
        print("photo")
        sys.exit(0)

    # Optional: load ONNX model from env (e.g. DOCUMENT_CLASSIFIER_MODEL=model.onnx)
    model_path = os.environ.get("DOCUMENT_CLASSIFIER_MODEL")
    if model_path and os.path.isfile(model_path):
        try:
            import onnxruntime as ort
            import numpy as np
            from PIL import Image
            sess = ort.InferenceSession(model_path)
            img = Image.open(image_path).convert("RGB")
            img = img.resize((224, 224))
            arr = np.array(img).astype(np.float32) / 255.0
            arr = arr.transpose(2, 0, 1)[np.newaxis, ...]
            out = sess.run(None, {sess.get_inputs()[0].name: arr})[0]
            # Assume output shape (1, 2) or (1,) with document=1, photo=0
            if out.shape[-1] >= 2:
                pred = "document" if out[0][1] > out[0][0] else "photo"
            else:
                pred = "document" if out.flat[0] > 0.5 else "photo"
            print(pred)
            return
        except Exception as e:
            sys.stderr.write(f"ONNX failed: {e}\n")
            # fallback to heuristic

    if looks_like_document_heuristic(image_path):
        print("document")
    else:
        print("photo")


if __name__ == "__main__":
    main()
