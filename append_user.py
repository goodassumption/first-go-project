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
    except requests.exceptions.RequestException as e:
        return None, f"Сетевая ошибка при обновлении счета: {e}"

    if response.ok:
        return response.json()
    else:
        return None, f"Ошибка сервера при обновлении счета: {response.status_code} {response.reason}"


def get_leaderboard(limit=5):
    url = f"{SERVER_URL}/leaderboard?limit={limit}"
    
    try:
        response = requests.get(url)
    except requests.exceptions.RequestException as e:
        return None, f"Сетевая ошибка при получении рейтинга: {e}"
    
    if response.ok:
        return response.json().get('entries', [])
    else:
        return None, f"Ошибка сервера при получении рейтинга: {response.status_code} {response.reason}"


def plus_count(user_id, user_name):
    return update_score(user_id, user_name, 1)

if __name__ == '__main__':
    
    try:
        new_user_id = "test_user_" + str(int(requests.get(SERVER_URL + "/leaderboard").elapsed.total_seconds() * 1000000))
    except requests.exceptions.RequestException:
        print("Ошибка: Не удалось подключиться к серверу.")
        exit()
        
    new_user_name = "{Name}"
    
    initial_user_resp = update_score(new_user_id, new_user_name, 0)
    
    user_inc = plus_count(new_user_id, new_user_name)
    
    leaderboard_data = get_leaderboard()
    
    if isinstance(leaderboard_data, list) and leaderboard_data:
        for entry in leaderboard_data:
            print(f"{entry['Rank']}. {entry['Name']} - {entry['Scr']}")
    else:
        print(f"Ошибка получения рейтинга: {leaderboard_data[1]}")