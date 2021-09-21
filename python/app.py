# python app.py

from flask import Flask
from flask import make_response
import requests

app = Flask(__name__)

@app.route('/hello')
def hello():
    return 'hello from python\n'

@app.route('/hello/proxy/<path:subpath>')
def proxy(subpath):
    parts = subpath.split('/')
    proxyTarget = parts[0]
    rest = '/'.join(parts[1:])

    url = f'http://{proxyTarget}/hello' if (len(parts) == 1) else f'http://{proxyTarget}/hello/proxy/{rest}'
    r = requests.get(url)
    return f'{r.text}hello from python\n'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=80)
