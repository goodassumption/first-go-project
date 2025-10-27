import requests
import json

SERVER_URL = "http://localhost:9090"

def update_score(id, name, scr):
    url = f"{SERVER_URL}/score/update"
    payload = {
        "id": id,
        "name": name,
        "addScr": scr
    }
    
    try:
        response = requests.post(url, json=payload)
        response.raise_for_status() 
        
        return response.json()
            
    except requests.exceptions.RequestException as e:
        return None, f"Ошибка при обновлении счета: {e}"

def get_leaderboard(limit=5):
    url = f"{SERVER_URL}/leaderboard?limit={limit}"
    
    try:
        response = requests.get(url)
        response.raise_for_status()
        
        return response.json().get('entries', [])
            
    except requests.exceptions.RequestException as e:
        return None, f"Ошибка при получении рейтинга: {e}"


def plus_count():
    return update_score(4, "user", 149)

if __name__ == '__main__':
    
    
    """initial_alice = update_score("1", "Alice", 100)
    initial_bob = update_score("2", "Bob", 150)
    initial_charlie = update_score("3", "Charlie", 50)"""
    
    leaderboard_data = get_leaderboard()
    
    id = str(len(leaderboard_data) + 1)
    initial_user = update_score(4, "user", 0)

    
    user_inc_1 = plus_count()
    
    user_inc_2 = plus_count()
    
    
    if isinstance(leaderboard_data, list) and leaderboard_data:
        for entry in leaderboard_data:
            print(f"{entry['Rank']}. {entry['Name']} {entry['Scr']}")
    elif isinstance(leaderboard_data, tuple) and leaderboard_data[0] is None:
        print(f"Ошибка получения рейтинга: {leaderboard_data[1]}")