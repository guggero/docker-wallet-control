$(function () {
    Handlebars.registerHelper('ago', function (options) {
        duration = moment.duration(moment().diff(moment(options.fn(this))));
        return duration.humanize();
    });
    reload();
});

function loader(cb) {
    $("#loader").show();
    return function (data) {
        $("#loader").hide();
        cb(data);
    }
}

function reload() {
    $.get('/summary', loader(function (data) {
        updateTemplate("#summary", "#content-placeholder", data);
        for (var i = 0; i < data.length; i++) {
            updateLogs(data[i].wallettype);
            updateContainerInfo(data[i].wallettype);
        }
    }));
}

function restart(wallet) {
    $("#loader").show();
    $.get('/' + wallet + '/restart' + avoidCache(), function (data) {
        updateContainerInfo(wallet);
    });
}

function updateTemplate(templateId, targetId, data) {
    var source = $(templateId).html();
    var template = Handlebars.compile(source);
    $(targetId).html(template(data));
}

function updateLogs(wallet) {
    $.get('/' + wallet + '/logs' + avoidCache(), loader(function (logs) {
        updateTemplate("#logs", "#logs-" + wallet, logs.reverse());
    }));
}

function updateContainerInfo(wallet) {
    $.get('/' + wallet + '/health' + avoidCache(), loader(function (health) {
        health.wallettype = wallet;
        updateTemplate("#container-info", "#container-info-" + wallet, health);
    }));
}

function createAccount(wallet) {
    var accountName = $("#create-account-" + wallet).val();
    $.get('/' + wallet + '/account/' + accountName, loader(function (data) {
        alert('Account ' + accountName + ' created, address: ' + data);
        reload();
    }));
}

function sendFrom(wallet) {
    var accountName = $("#send-account-" + wallet).val();
    var toAddress = $("#send-toaddress-" + wallet).val();
    var amount = $("#send-amount-" + wallet).val();
    var post = {
        account: accountName,
        address: toAddress,
        amount: amount * 1
    };
    $.ajax({
        type: 'POST',
        url: '/' + wallet + '/sendfrom',
        data: JSON.stringify(post),
        contentType: 'application/json',
        success: loader(function (result) {
            alert('Sent ' + amount + ' to ' + toAddress + ' from account ' + accountName + ', result: ' + result);
            reload();
        }),
        dataType: 'json'
    });
}

function avoidCache() {
    return '?_=' + (new Date().getTime());
}
