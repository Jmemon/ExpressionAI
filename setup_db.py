import sqlite3


def create_db():
    conn = sqlite3.connect('dialogue.db')
    c = conn.cursor()
    c.execute('''CREATE TABLE IF NOT EXISTS dialogues
                 (id INTEGER PRIMARY KEY, client_text TEXT, model_response TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)''')
    conn.commit()
    conn.close()
    

if __name__ == "__main__":
    create_db()
    print("Database setup completed.")
