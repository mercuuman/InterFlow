const config = {
  baseUrl: "http://localhost:8080",
  headers: {
    'Content-Type': 'application/json'
  }
}

function checkStatus(res) {
  return res.ok?res.json():Promise.reject(`Ошибка: ${res.status}`);
}

export const registrationUser = (requestConfig) => {
  return fetch(`${config.baseUrl}/api/registration`, {
    method: 'POST',
    headers: config.headers,
    body: JSON.stringify({
      name: requestConfig.name,
      mail: requestConfig.mail,
      password: requestConfig.password
    })
  })
  .then(res => {
    return checkStatus(res);
  })
}


