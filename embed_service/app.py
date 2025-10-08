from flask import Flask, request, jsonify
from sentence_transformers import SentenceTransformer

app = Flask(__name__)
model = SentenceTransformer('sentence-transformers/all-MiniLM-L6-v2')


@app.route('/embed', methods=['POST'])
def embed():
    data = request.get_json(force=True, silent=True) or {}
    texts = data.get('texts', [])
    if not isinstance(texts, list) or not texts:
        return jsonify({'error': 'texts must be a non-empty list'}), 400
    vectors = model.encode(texts, normalize_embeddings=True).tolist()
    return jsonify({'vectors': vectors})


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000)


