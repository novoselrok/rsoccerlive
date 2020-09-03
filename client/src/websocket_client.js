const RECONNECT_TIMEOUT_MS = 5000

export default class WebSocketClient {
  constructor (wsEndpoint, onEventFn) {
    this.wsEndpoint = wsEndpoint
    this.onEventFn = onEventFn
    this.ws = null

    this._connect = this._connect.bind(this)
    this._reconnect = this._reconnect.bind(this)
    this._onClose = this._onClose.bind(this)
    this._onMessage = this._onMessage.bind(this)
    this._onError = this._onError.bind(this)

    this._connect()
  }

  _connect () {
    this.ws = new WebSocket(this.wsEndpoint)
    this.ws.addEventListener('message', this._onMessage)
    this.ws.addEventListener('error', this._onError)
    this.ws.addEventListener('close', this._onClose)
  }

  _reconnect() {
    // eslint-disable-next-line
    console.log('WebSocket reconnecting')
    this.ws.removeEventListener('message', this._onMessage)
    this.ws.removeEventListener('error', this._onError)
    this.ws.removeEventListener('close', this._onClose)
    this._connect()
  }

  _onClose () {
    const randomReconnectTimeoutMs = RECONNECT_TIMEOUT_MS + Math.floor(Math.random() * RECONNECT_TIMEOUT_MS)
    setTimeout(this._reconnect, randomReconnectTimeoutMs)
  }

  _onMessage (e) {
    this.onEventFn(JSON.parse(e.data))
  }

  _onError (e) {
    // eslint-disable-next-line
    console.error('WebSocket error', e)
    this.ws.close()
  }
}
