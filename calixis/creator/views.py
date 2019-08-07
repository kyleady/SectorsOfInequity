from django.http import JsonResponse
import requests

def index(request):
    r = requests.post('http://127.0.0.1:8080/grid/', json={})
    output = r.json()
    return JsonResponse(output)
