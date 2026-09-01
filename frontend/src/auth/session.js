export function getAccessToken() {
  return localStorage.getItem('infrapilot.accessToken');
}

export function setAccessToken(token) {
  localStorage.setItem('infrapilot.accessToken', token);
}

export function clearSession() {
  localStorage.removeItem('infrapilot.accessToken');
}
