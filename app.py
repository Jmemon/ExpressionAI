
from flask import Flask, request, jsonify
import sqlite3
from langchain.llms import HuggingFace


app = Flask(__name__)
lm = HuggingFace(model_name="EleutherAI/gpt-neo-2.7B")


def insert_dialogue(client_text, model_response):
    conn = sqlite3.connect('dialogue.db')
    c = conn.cursor()
    c.execute("INSERT INTO dialogues (client_text, model_response) VALUES (?, ?)", (client_text, model_response))
    conn.commit()
    conn.close()


@app.route('/process_text', methods=['POST'])
def process_text():
    data = request.get_json()
    input_text = data.get('text', '')
    
    if not input_text:
        return jsonify({"error": "No text provided"}), 400

    try:
        response = lm.generate(input_text)
        insert_dialogue(input_text, response)
        return jsonify({"response": response})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


if __name__ == '__main__':
    app.run(debug=True, port=5000)
