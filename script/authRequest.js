document.querySelector('.signButton').addEventListener('click', function() {
    let login = document.querySelector('input[name="login"]').value
    let password = document.querySelector('input[name="password"]').value

    fetch('http://localhost:8080/auth', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            login: login,
            password: password
        })
    })
    .then(response => response.json())
    .then(data => {
        alert(data.result)
    })
    .catch(error => {
        alert('Error: ' + error)
    })
})

document.querySelector('.showTableButton').addEventListener('click', function(){
    let DBTable = document.querySelector('.DBTable')
    DBTable.querySelectorAll('div:not([class])').forEach(elem => elem.remove())
    DBTable.style.opacity='1'

    function createCell(text){
        let cell = document.createElement('div')
        cell.innerHTML=text
        cell.style.fontFamily='sans-serif'
        cell.style.padding="5px 5px 5px 5px"
        cell.style.backgroundColor="white"
        cell.style.color="blue"
        cell.style.border="1px solid blue"
        DBTable.appendChild(cell)
    }

    fetch('http://localhost:8080/show', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        },
    })
    .then(response => response.json())
    .then(data => {
        data.forEach(elem => {
            createCell(elem.login)
            createCell(elem.password)
        })
    })
    .catch(error => {
        alert('Не вышло отобразить базу данных' + error)
    })
})