$('.signButton').click(function() {
    let login = $('input[name="login"]').val()
    let password = $('input[name="password"]').val()

    $.ajax({
        url: 'http://localhost:8080/auth',
        type: 'POST',
        data: JSON.stringify({
            login: login,
            password: password
        }),
        contentType: 'application/json'
    })
    .done(function(data) {
        alert(data.result)
    })
    .fail(function(xhr, status, error) {
        alert("Error: " + error)
    })
})
