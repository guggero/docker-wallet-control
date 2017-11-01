angular.element(document.body).ready(function () {
  angular.bootstrap(document.body, ['app'])
});

function avoidCache() {
  return '?_=' + (new Date().getTime());
}

function formatString(str) {
  var args = [].slice.call(arguments, 1),
    i = 0;

  return str.replace(/%s/g, function () {
    return args[i++];
  });
}

angular
  .module('app', [])
  .constant('moment', window.moment)
  .constant('lodash', window._)
  .filter('ago', function (moment) {
    return function (input) {
      var duration = moment.duration(moment().diff(moment(input)));
      return duration.humanize()
    }
  })
  .filter('amount', function () {
    return function (input) {
      return (input / 1e8).toFixed(4).replace(/(\d)(?=(\d{3})+\.)/g, '$1\'');
    }
  })
  .filter('date', function (moment) {
    return function (input) {
      return moment(input).format('YYYY-MM-DD hh:mm:ss')
    }
  })
  .filter('change', function (lodash) {
    return function (tx, daysBack) {
      var change = 0;
      var now = moment();
      var days = daysBack || 999999;
      lodash.forEach(tx, function (t) {
        var then = moment(t.time_utc);
        if (now.diff(then, 'days') <= days) {
          change += t.change;
        }
      });
      return change;
    }
  })
  .filter('formatString', function () {
    return formatString;
  })
  .component('app', {
    templateUrl: 'app.html',
    controller: AppController,
    controllerAs: 'vm',
    bindings: {}
  });


function AppController($http, $q, $scope) {
  var vm = this;

  vm.wallets = [];
  vm.uiData = null;
  vm.tab = null;
  vm.running = true;
  vm.masternodes = [];

  vm.activate = activate;
  vm.changeTab = changeTab;
  vm.getLogs = getLogs;
  vm.getContainerInfo = getContainerInfo;
  vm.restart = restart;
  vm.createAccount = createAccount;
  vm.sendFrom = sendFrom;
  vm.sendCommand = sendCommand;

  activate();

  ////////////////

  function activate() {
    vm.running = true;
    vm.masternodes = [];
    return $http.get('/summary' + avoidCache()).then(function (response) {
      vm.wallets = response.data.summaries;
      vm.uiData = response.data.uiData;

      if (vm.tab === null && vm.wallets.length > 0) {
        vm.tab = vm.wallets[0].label;
      }

      var promises = [];
      vm.wallets.forEach(function (wallet) {
        promises.push(getContainerInfo(wallet));
        promises.push(getLogs(wallet));
        promises.push(getMasternodeStats(wallet));
      });
      $q.all(promises).finally(function () {
        vm.running = false;
      })
    })
  }

  function changeTab(wallet) {
    vm.tab = wallet.label;
  }

  function getContainerInfo(wallet) {
    return $http.get('/' + wallet.containername + '/health' + avoidCache()).then(function (response) {
      wallet.container = response.data;
    });
  }

  function getLogs(wallet) {
    return $http.get('/' + wallet.containername + '/logs' + avoidCache()).then(function (response) {
      wallet.logs = response.data.reverse().join('\n');
    })
  }

  function getMasternodeStats(wallet) {
    if (wallet.masternodeStatus && wallet.masternodeStatus.addr) {
      var mnStatus = wallet.masternodeStatus;
      mnStatus.service = mnStatus.service || mnStatus.netaddr;
      mnStatus.status = mnStatus.message || mnStatus.status;
      mnStatus.address = mnStatus.payee || mnStatus.addr;

      var url = formatString(vm.uiData.apis.address, wallet.wallettype, wallet.masternodeStatus.addr);
      return $http.get(url).then(function (response) {
        var mnStatus = wallet.masternodeStatus;
        vm.masternodes.push({
          service: mnStatus.service,
          status: mnStatus.status,
          type: wallet.wallettype,
          address: mnStatus.address,
          balance: response.data.addresses[0].final_balance,
          transactions: mergeTransactions(response.data.txs, 'hash', 'change')
        });
      });
    } else {
      return $q.when();
    }
  }

  function restart(wallet) {
    vm.running = true;
    $http.get('/' + wallet.containername + '/restart' + avoidCache()).then(function () {
      return getContainerInfo(wallet);
    }).finally(function () {
      vm.running = false;
    });
  }

  function createAccount(wallet) {
    vm.running = true;
    var accountName = $('#create-account-' + wallet.containername).val();
    $http.get('/' + wallet.containername + '/account/' + accountName).then(function (response) {
      alert('Account ' + accountName + ' created, address: ' + response.data);
      activate();
    });
  }

  function sendFrom(wallet) {
    var accountName = $('#send-account-' + wallet.containername).val();
    var toAddress = $('#send-toaddress-' + wallet.containername).val();
    var amount = $('#send-amount-' + wallet.containername).val();
    var post = {
      account: accountName,
      address: toAddress,
      amount: amount * 1
    };
    $http.post('/' + wallet.containername + '/sendfrom', JSON.stringify(post)).then(function (response) {
      alert('Sent ' + amount + ' to ' + toAddress + ' from account ' + accountName + ', result: ' + response.data);
      activate();
    });
  }

  function sendCommand(wallet) {
    var commandParts = $('#send-comand-' + wallet.containername).val().split(' ');
    var post = {
      method: commandParts.shift(),
      args: commandParts
    };
    $http.post('/' + wallet.containername + '/command', JSON.stringify(post)).then(function (response) {
      if (typeof response.data === 'string') {
        wallet.commandResponse = response.data;
      } else {
        wallet.commandResponse = JSON.stringify(response.data, null, 2);
      }
    });
  }

  function mergeTransactions(tx, key, mergeProperty) {
    for (var i = 0; i < tx.length; i++) {
      if (i < tx.length - 1) {
        var a = tx[i];
        while (i + 1 < tx.length && a[key] === tx[i + 1][key]) {
          a[mergeProperty] += tx[i + 1][mergeProperty];
          tx.splice(i + 1, 1);
        }
      }
    }
    return tx;
  }
}
