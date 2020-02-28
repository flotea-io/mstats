var api_ip = "http://0.0.0.0:8099";
var api_login = "http://0.0.0.0:9940";
var MachineAPP = (function () {
    function MachineAPP() {
        var _this = this;
        this.minersNamesArray = [];
        this.all_mhs = 0;
        this.all_temp = 0;
        this.machine_quantity = 0;
        this.numberOfDeadMachines = 0;
        this.numberOfOverhiting = 0;
        this.numberOfdeadCards = 0;
        this.listOfdeadCards = "";
        this.sortType = "name";
        this.sortAliveData = [];
        this.sortDeadData = [];
        this.pins = [];
        this.addedPins = [];
        this.pinsSetting = [];
        this.revert = true;
        this.changeWatch = false;
        this.walletsList = [];
        this.configsList = [];
        if (MachineAPP._instance) {
            throw new Error("Error: Instantiation failed: Use MachineAPP.getInstance() instead of new.");
        }
        MachineAPP._instance = this;
        this.minersArray = [];
        $(document).ready(function () { return _this.initHandlers(); });
    }
    MachineAPP.prototype.initHandlers = function () {
        var _this = this;
        //this.checkingIsEmpty();
        this.isLogged();
        this.fillConfigList();
        this.fillWalletsList();
        this.getMiners();
        $(document).on("click", "#navbarCollapse [data-target='#addMachine']", function () { return _this.addMinersList(); });
        $(document).on("click", "#logInButton", function () { return _this.logIn(); });
        $(document).on("click", ".reboot", function (event) { return _this.menageMachine(event); });
        $(document).on("click", ".restart", function (event) { return _this.menageMachine(event); });
        $(document).on("click", ".delete", function (event) { return _this.menageMachine(event); });
        $(document).on("click", ".on_off", function (event) { return _this.menageMachine(event); });
        $(document).on("click", ".hard-reset", function (event) { return _this.menageMachine(event); });
        $(document).on("click", ".editHardwareButton", function () { return _this.editHardware(); });
        $(document).on("click", ".wallet_menagment", function () { return _this.fill_wallet_menagment(); });
        $(document).on("click", ".config_menagment", function () { return _this.fill_config_menagment(); });
        $(document).on("click", "#changePinsNumber", function () { return _this.changePinsNumber(); });
        $(document).on("click", ".wallet_edit", function (event) { return _this.editWallet(event); });
        $(document).on("click", ".wallet_add", function () { return _this.addWallet(); });
        $(document).on("click", ".wallet_delete", function (event) { return _this.deleteWallet(event); });
        $(document).on("click", ".wallet_isDefault", function (event) { return _this.isDefaultWallet(event); });
        $(document).on("click", ".config_edit", function (event) { return _this.editConfig(event); });
        $(document).on("click", ".config_add", function () { return _this.addConfig(); });
        $(document).on("click", ".config_delete", function () { return _this.deleteConfig(); });
        $("#addEditWallet").find(".btn-outline-success").on("click", function () { return _this.walletAction(); });
        $("#addEditConfig").find(".btn-outline-success").on("click", function () { return _this.configAction(); });
        $(document).on("click", ".menage_machine", function () { return _this.changeWallet(); });
        //  $(document).on("click", ".edit", (event: any) => this.editMachine(event));
        //  $(document).on("click", ".machine_operation .toConsol", (event: any) => this.toConsole(event));
        $('#statistics').on('shown.bs.modal', function () { return _this.statistics(); });
        this.sortHandlers();
        $(".filter").on("keyup", function () { return _this.filter(); });
        this.loadPinsSettings(function () { return _this.enableHardButons(); });
        $("#editHardwareModal").on("click", "#clearPins", function () {
            _this.pinsSetting = [];
            _this.generatePins();
            $.ajax({
                type: "POST",
                url: api_ip + "/arduino/pins/clear",
                header: ({ authentication: _this.cookie.access_token }),
                crossDomain: true,
                dataType: "json",
                contentType: "application/json"
            });
            _this.renderMachinesResetOptions();
            _this.renderPins();
        });
        this.refreshTemperatureDate();
        setInterval(function () { return _this.refreshResetTime(); }, 1000);
    };
    MachineAPP.getInstance = function () {
        return MachineAPP._instance;
    };
    MachineAPP.prototype.isLogged = function () {
        this.cookie = JSON.parse(Cookies.get('authentication'));
        $.ajax({
            url: api_login + "/rest/logged",
            type: "get",
            contentType: "application/json",
            Authorization: "Bearer " + this.cookie.access_token,
            dataType: 'json'
        }).done(function (respond) {
            console.log("respond: ", respond);
        });
    };
    MachineAPP.prototype.logIn = function () {
        var _this = this;
        var nick = $("#logInNick").val();
        var password = $("#logInPassword").val();
        //let remember = $("#logInRemember").prop("checked");
        $.ajax({
            type: "POST",
            url: api_login + "/oauth2/token",
            crossDomain: true,
            dataType: "json",
            contentType: "application/json",
            data: JSON.stringify({
                grant_type: "password",
                username: nick,
                password: password,
                client_id: "testclient",
                client_secret: "testpass"
            })
        }).done(function (respond) {
            Cookies.set('authentication', JSON.stringify(respond), { expiry: respond.expires_in });
            _this.cookie = JSON.parse(Cookies.get('authentication'));
        });
    };
    MachineAPP.prototype.filter = function () {
        var written_machine = $(".filter:focus").val();
        $('.filter').each(function () {
            $('.filter').val(written_machine);
        });
        if (written_machine != "") {
            for (var _i = 0, _a = this.minersNamesArray; _i < _a.length; _i++) {
                var name_1 = _a[_i];
                if (name_1.indexOf(written_machine) == -1) {
                    $("#machine_" + name_1).parent().hide();
                }
                else
                    $("#machine_" + name_1).parent().show();
            }
        }
        else {
            for (var _b = 0, _c = this.minersNamesArray; _b < _c.length; _b++) {
                var name_2 = _c[_b];
                $("#machine_" + name_2).parent().show();
            }
        }
    };
    MachineAPP.prototype.sortHandlers = function () {
        var _this = this;
        $(".sortByName").on("click", function () {
            $(".sortByDropdown .active").removeClass("active");
            $(".sortByName").addClass("active");
            _this.sortType = "name";
            _this.sortBy(_this.sortType);
        });
        $(".sortByRunning").on("click", function () {
            $(".sortByDropdown .active").removeClass("active");
            $(".sortByRunning").addClass("active");
            _this.sortType = "runTime";
            _this.sortBy(_this.sortType);
        });
        $(".sortByMh").on("click", function () {
            $(".sortByDropdown .active").removeClass("active");
            $(".sortByMh").addClass("active");
            _this.sortType = "mh";
            _this.sortBy(_this.sortType);
        });
        $(".sortByTemp").on("click", function () {
            $(".sortByDropdown .active").removeClass("active");
            $(".sortByTemp").addClass("active");
            _this.sortType = "heat";
            _this.sortBy(_this.sortType);
        });
    };
    MachineAPP.prototype.sortBy = function (type) {
        var asc = true;
        if (type == "runTime" || type == "mh" || type == "heat") {
            asc = false;
        }
        this.sortAliveData.sort(function (a, b) {
            if (a[type] < b[type])
                return asc ? -1 : 1;
            if (a[type] > b[type])
                return asc ? 1 : -1;
            return 0;
        });
        this.sortDeadData.sort(function (a, b) {
            if (a[name] < b[name])
                return asc ? -1 : 1;
            if (a[name] > b[name])
                return asc ? 1 : -1;
            return 0;
        });
        var container = $("#machine_container");
        var deadContainer = $("#dead_machine_container");
        for (var _i = 0, _a = this.sortAliveData; _i < _a.length; _i++) {
            var m = _a[_i];
            container.append($("#machine_" + m.name).parent());
        }
        for (var _b = 0, _c = this.sortDeadData; _b < _c.length; _b++) {
            var m = _c[_b];
            deadContainer.append($("#machine_" + m.name).parent());
        }
    };
    MachineAPP.prototype.addMinersList = function () {
        var _this = this;
        $.ajax({
            url: api_ip + "/clients/disconected",
            type: "get",
            dataType: 'json'
        }).done(function (minersToAdd) {
            $("#addMachine").removeClass("d-none");
            $("#addMachine .modal-body").text(" ");
            var addMachineContainer = $(".addMachineList").clone();
            for (var _i = 0, minersToAdd_1 = minersToAdd; _i < minersToAdd_1.length; _i++) {
                var i = minersToAdd_1[_i];
                if (i != "") {
                    addMachineContainer.find(".machinesToAdd").first().attr("id", "machineToAdd_" + i);
                    addMachineContainer.find(".nameToAdd").text(i);
                    addMachineContainer.find(".machinesToAdd").first().removeClass("d-none").addClass("d-flex");
                    $("#addMachine .modal-body").append(addMachineContainer.html());
                }
            }
            _this.addMachine();
            if (minersToAdd[0] == "" || minersToAdd.length == 0) {
                $("#addMachine .modal-body").text("No machines to add");
            }
        }).fail(function () {
            $("#addMachine").addClass("d-none");
        });
    };
    MachineAPP.prototype.refreshResetTime = function () {
        if (typeof this.temperatureData == "undefined")
            return;
        // setTimeout(() => {
        //   this.refreshMinersData();
        // }, 5000)
    };
    MachineAPP.prototype.getMiners = function () {
        var _this = this;
        $.ajax({
            url: api_ip + "/clients",
            type: "get",
            dataType: 'json'
        }).done(function (miners) {
            if (!$("#alert-infos").hasClass("d-none")) {
                $("#alert-infos").addClass("d-none");
            }
            $("#dead_machine_container").text("");
            $("#dead_machine_container").addClass("d-none");
            $("#machine_container").text("");
            _this.minersArray = miners;
            _this.renderMiners();
            setInterval(function () {
                _this.refreshMinersData();
            }, 12000);
            _this.refreshMinersData();
        }).fail(function () {
            $("#alert-infos").removeClass("d-none");
            $("#alert-infos").text("Error when try connect with machines");
        });
    };
    MachineAPP.prototype.renderMiners = function () {
        for (var _i = 0, _a = Object.keys(this.minersArray); _i < _a.length; _i++) {
            var miner_key = _a[_i];
            var machine = $('#template_machine>').clone();
            machine.find(">div").addClass("miner");
            machine.find(">div").first().attr("id", "machine_" + miner_key);
            this.minersNamesArray.push(miner_key);
            machine.find(">div").first().attr("data-id", miner_key);
            machine.find(".miner_name").text(this.minersArray[miner_key].MinerName);
            $("#machine_container").append(machine);
        }
        ;
        //Function needed to show popups correctly
        $(".editHardwareButton").addClass("d-md-flex");
    };
    MachineAPP.prototype.refreshMinersData = function () {
        var _this = this;
        $.ajax({
            url: api_ip + "/latest/stat",
            type: "get",
            dataType: 'json'
        }).done(function (respond) {
            _this.all_mhs = 0;
            _this.all_temp = 0;
            _this.machine_quantity = 0;
            _this.numberOfDeadMachines = 0; //Refresh errors [END]
            $("#showIssues .overhiting").text(""); //Refresh errors [START]
            $("#showIssues .deadCards").text("");
            $(".issuesButton").addClass("d-none");
            $("#overhiting_counter").text("");
            $("#deadCards_counter").text("");
            $("#overhiting_span").addClass("d-none");
            $("#deadCards_span").addClass("d-none");
            $("#deadMachines_span").addClass("d-none");
            _this.overhiting_helper = " ";
            _this.deadCards_helper = " ";
            _this.listOfdeadCards = " ";
            _this.numberOfOverhiting = 0;
            _this.numberOfdeadCards = 0;
            _this.sortAliveData = [];
            _this.sortDeadData = [];
            for (var _i = 0, _a = Object.keys(_this.minersArray); _i < _a.length; _i++) {
                var miner_key = _a[_i];
                var miner = _this.minersArray[miner_key];
                var gpus = respond.hardware[miner_key];
                gpus = typeof gpus == "undefined" ? gpus = null : gpus = JSON.parse(respond.hardware[miner_key].gpuList);
                if (respond.miners[miner.MinerName] === undefined || gpus === null) {
                    _this.deadMiners("machine_" + miner_key);
                    continue;
                }
                _this.fillMiners("machine_" + miner_key, respond.miners[miner.MinerName], gpus, respond.config[miner_key]);
            }
            if (_this.changeWatch == true) {
                _this.sortBy(_this.sortType);
            }
            _this.all_mhs = _this.Round(_this.all_mhs, 0); //Refreshing info of total Mh/s and Temp on menu [START]
            _this.all_temp = _this.Round(_this.all_temp / _this.machine_quantity, 0);
            $(document).find(".total_data").html("<span class='d-block d-md-inline-block'>Total:  " + _this.all_mhs + "Mh/s | "
                + "<span data-tooltip='yes' title='Averge temperature from all machines'>" + _this.all_temp + "&#186C | </span></span>"
                + "<span data-tooltip='yes' title='Working machines'> Machines: <img src='style/icons/icon_machine_works.svg' alt='Working machines' width='20' height='15'/> "
                + _this.machine_quantity + "</span>"
                + "<span data-tooltip='yes' title='Not working machines' > <img src='style/icons/icon_machine_not_works.svg' alt='Not working machines' width='20' height='15'/> "
                + _this.numberOfDeadMachines) + "</span>"; //[END]
            if ((_this.numberOfOverhiting + _this.numberOfdeadCards + _this.numberOfDeadMachines) > 0) {
                //If no errors, in modal we can see comunicate 'No errors'
                $(".issuesButton").removeClass("d-none");
                $(".numberOfIssues").text("Issues(" + (_this.numberOfOverhiting + _this.numberOfdeadCards + _this.numberOfDeadMachines) + ")");
            }
            $("[data-tooltip='yes']").tooltip({ boundary: 'window' });
            _this.changeWatch = false;
        }).fail(function () {
            console.log("Error when try refresh");
        });
    };
    MachineAPP.prototype.Round = function (n, k) {
        var factor = Math.pow(10, k);
        return Math.round(n * factor) / factor;
    };
    MachineAPP.prototype.measureTime = function (minutes) {
        var day = moment.duration(minutes, "minutes").humanize();
        var hours = moment.duration(minutes / 60, "minutes").humanize();
        var seconds = moment.duration(minutes / 120, "minutes").humanize();
        var day_test = /day/i;
        if (day_test.test(day)) {
            return day + " and " + hours;
        }
        else {
            return day + " and " + seconds;
        }
    };
    MachineAPP.prototype.fillMiners = function (id, miner, fans, wallet_info) {
        var _this = this;
        this.machine_quantity++;
        if ((moment() - moment(miner.Time)) > 90000) {
            this.deadMiners(id); //Checking if log is not to old, if is - machine is down
            this.machine_quantity--;
            $("#" + id).find(".miner_issues").text("Last log of machine is from " + moment(miner.Time, "YYYYMMDDHHmmss").fromNow());
            return;
        }
        else if ($("#" + id + " .card-body>div").first().hasClass("dead_machine_container")) {
            this.changeWatch = true;
            $("#" + id + " .card-body").html($("#template_machine .card-body>").clone());
            $("#machine_container").append($("#" + id).parent());
            $("#dead_machine_container").find("#" + id).parent().remove();
        }
        $("#" + id).find(".miner_issues").text("");
        if (wallet_info.Currency == "ETH") {
            $("#" + id).find(".currency").html('<img src="style/icons/icon_eth.svg" alt="ETH Graphic"/>' + wallet_info.Currency);
        }
        if (wallet_info.Currency == "ETC") {
            $("#" + id).find(".currency").html('<img src="style/icons/icon_etc.svg" alt="ETC Graphic"/>' + wallet_info.Currency);
        }
        $("#" + id).find(".wallet").text(" " + wallet_info.WalletName);
        var machine_name = id.replace("machine_", "");
        $("#" + id).find(".miner_name").text(machine_name);
        var arrayMghs = miner.DetailedEthHashRatePerGPU.split(';');
        var arrayTemp = miner.Temperatures.split(';');
        var forTemp = 0;
        var avergeTemp = 0;
        var oneGpu;
        var maxTemp = 0;
        for (var i = 0; i < arrayMghs.length; i++) {
            if (!$("#" + id + " #gpus_container").find(".row").hasClass("card-" + i)) {
                $("#" + id + " #gpus_container").append($('#gpus_template_element>div').clone().addClass("card-" + i));
            }
            oneGpu = $("#" + id + " .card-" + i);
            oneGpu.find(".gpu_name").text("GPU" + i);
            oneGpu.find(".gpu_mhs").text(this.Round((arrayMghs[i] / 1024), 1));
            oneGpu.find(".gpu_temp").html(arrayTemp[forTemp] + " &#186C");
            avergeTemp += parseInt(arrayTemp[forTemp]);
            if (maxTemp < parseInt(arrayTemp[forTemp])) {
                maxTemp = parseInt(arrayTemp[forTemp]);
            }
            this.checkingErrors(oneGpu, arrayTemp[forTemp], id, "GPU" + (i), this.Round((arrayMghs[i] / 1024), 1));
            oneGpu.find(".switch_change_id").attr("id", id + "_switch_" + i);
            oneGpu.find(".switch_label").attr("for", id + "_switch_" + i);
            oneGpu.find(".switch_change_id_mobile").attr("id", id + "_switch_mobile_" + i);
            oneGpu.find(".switch_label_mobile").attr("for", id + "_switch_mobile_" + i);
            $("#" + id + " .card-" + i).find(".fan_select").addClass("fan_select_" + machine_name + "_" + i)
                .attr({ "data-miner": machine_name, "data-gpu": i + 1 })
                .val(fans[i + 1].declaredFanSpeed).change(function (event) { return _this.fanSpeedChanged(event); });
            oneGpu.find(".fan_speed").addClass("fan_speed_" + machine_name + "_" + i);
            oneGpu.find(".fan_speed_" + machine_name + "_" + i).text(fans[i + 1].fanSpeed);
            forTemp += 2;
        }
        $("#" + id).find(".collapseButton").click(function (event) {
            if (!$(event.currentTarget).parents(".mobile").find(".h-0").hasClass("show")) {
                $(event.currentTarget).parents(".gpus_container").find(".collapseArrow").css('transform', 'rotate(' + 360 + 'deg)');
                ;
                $(event.currentTarget).parents(".gpus_container").find(".show").removeClass("show");
                $(event.currentTarget).parents(".mobile").find(".h-0").addClass("show");
                $(event.currentTarget).parents(".mobile").find(".collapseArrow").css('transform', 'rotate(' + 180 + 'deg)');
            }
            else {
                $(event.currentTarget).parents(".mobile").find(".h-0").removeClass("show");
                $(event.currentTarget).parents(".mobile").find(".collapseArrow").css('transform', 'rotate(' + 360 + 'deg)');
                ;
            }
        });
        avergeTemp = this.Round((avergeTemp / arrayMghs.length), 0);
        var totalMhs = this.Round((miner.TotalEthHashRate / 1024), 1);
        var totalShares = miner.EthShares;
        var runningTime = this.measureTime(miner.RunningTime);
        this.all_mhs += totalMhs;
        this.all_temp += avergeTemp;
        $("#" + id).find(".miner_total_shares").text(totalShares);
        $("#" + id).find(".miner_total_mhs").text(totalMhs);
        $("#" + id).find(".miner_total_temp").html(maxTemp + " &#186C");
        $("#" + id).find(".miner_total_time").text("Start running from: " + runningTime);
        $("#overhiting_counter").text("(" + this.numberOfOverhiting + ")");
        $("#deadCards_counter").text("(" + this.numberOfdeadCards + ")");
        this.enableHardButons();
        this.sortAliveData.push({ "name": machine_name, "heat": maxTemp, "mh": totalMhs, "runTime": miner.RunningTime });
    };
    MachineAPP.prototype.enableHardButons = function () {
        $(".hard-reset, .on_off").addClass("d-none");
        for (var i = 0; i < this.pinsSetting.length; i++) {
            if (this.pinsSetting[i].Function != null) {
                $("#machine_" + this.pinsSetting[i].MinerName + " ." + ((this.pinsSetting[i].Function == 0) ? "hard-reset" : "on_off")).removeClass("d-none");
            }
        }
    };
    MachineAPP.prototype.fanSpeedChanged = function (event) {
        $.ajax({
            type: "POST",
            url: api_ip + "/fans",
            header: ({ authentication: this.cookie.access_token }),
            crossDomain: true,
            dataType: "json",
            contentType: "application/json",
            data: JSON.stringify({
                id: event.target.dataset.gpu,
                machine: event.target.dataset.miner,
                speed: event.target.value
            })
        });
    };
    MachineAPP.prototype.deadMiners = function (id) {
        $("#deadMachines_span").removeClass("d-none");
        $("#showIssues .deadMachines").removeClass("d-none");
        var deadMinersTemplete = $('#template_dead_machine').clone().html();
        var machine_name = id.replace("machine_", "");
        this.sortDeadData.push({ "name": machine_name, "heat": 0, "mh": 0, "runTime": 0 });
        if (!$("#" + id + " .card-body>div").first().hasClass("dead_machine_container")) {
            this.changeWatch = true;
            $("#" + id + " .card-body").html(deadMinersTemplete);
            $("#dead_machine_container").append($("#" + id).parent());
            $("#dead_machine_container").removeClass("d-none");
            $("#machine_container").find("#" + id).parent().remove();
        }
        $("#" + id).find(".miner_name").text(machine_name);
        $("#" + id).find(".miner_issues").text("No information about machine in database");
        this.numberOfDeadMachines++;
        this.listOfdeadCards += machine_name + " ";
        $("#deadMachines_counter").text("(" + this.numberOfDeadMachines + ")");
        $(".deadMachines").text(this.listOfdeadCards);
    };
    MachineAPP.prototype.checkingErrors = function (oneGpu, arrayTemp, id, card, card_value) {
        $(oneGpu).find(".card-status").removeClass("bg-primary bg-dark bg-success bg-warning bg-danger text-dark");
        $(oneGpu).find(".row").removeClass("alert-dark text-dark");
        $(oneGpu).find(".gpu_info").text("");
        var machine_number = id.replace("machine_", "");
        if (arrayTemp == 0) {
            $(oneGpu).find(".card-status").addClass("bg-secondary");
            this.deadCards_stat(oneGpu, machine_number, card);
        }
        if (arrayTemp > 0 && arrayTemp <= 50)
            $(oneGpu).find(".card-status").addClass("bg-primary");
        if (arrayTemp > 50 && arrayTemp <= 65)
            $(oneGpu).find(".card-status").addClass("bg-success");
        if (arrayTemp > 65 && arrayTemp <= 75)
            $(oneGpu).find(".card-status").addClass("bg-warning");
        if (arrayTemp > 75 && arrayTemp <= 85) {
            $(oneGpu).find(".card-status").addClass("bg-danger");
            $(oneGpu).addClass("alert-danger");
            this.overhiting_stat(oneGpu, machine_number, card);
        }
        if (arrayTemp > 85) {
            $(oneGpu).find(".card-status").addClass("bg-danger");
            $(oneGpu).addClass("alert-danger");
            this.overhiting_stat(oneGpu, machine_number, card);
        }
        if (card_value == 0) {
            this.deadCards_stat(oneGpu, machine_number, card);
        }
    };
    MachineAPP.prototype.overhiting_stat = function (oneGpu, machine_number, card) {
        $("#overhiting_span").removeClass("d-none");
        $("#showIssues .overhiting").removeClass("d-none");
        $(oneGpu).find(".gpu_info").html("<img src='style/icons/icon_fire.svg' alt='Fire!' width='15px' data-tooltip='yes' title='Card is burning!'/>");
        if (this.overhiting_helper != machine_number) {
            var nav_item_overhiting = $("<div>")
                .text(machine_number + ": ")
                .attr("href", "#machine_" + machine_number)
                .attr("id", "machine_" + machine_number + "_overhiting");
            $("#showIssues .overhiting").append(nav_item_overhiting);
            $("#machine_" + machine_number + "_overhiting").append(" " + card);
            this.overhiting_helper = machine_number;
            this.numberOfOverhiting++;
        }
        else {
            $("#machine_" + machine_number + "_overhiting").append(", " + card);
            this.numberOfOverhiting++;
        }
    };
    MachineAPP.prototype.deadCards_stat = function (oneGpu, machine_number, card) {
        $(oneGpu).find(".row").removeClass("alert-primary alert-dark alert-success alert-warning alert-danger bg-danger bg-secondary text-dark");
        $(oneGpu).find(".row").addClass("alert-dark text-dark");
        $(oneGpu).find(".gpu_info").html("<img src='style/icons/icon_dead.svg' alt='Dead card!' width='20px' data-tooltip='yes' title='Card disabled'/>");
        $("#showIssues #deadCards_span").removeClass("d-none");
        $("#showIssues .deadCards").removeClass("d-none");
        if (!$(document).find(".deadCards").hasClass("d-none")) {
            if (this.deadCards_helper != machine_number) {
                var nav_item_deadCards = $("<div>")
                    .text(machine_number + ": ")
                    .attr("href", "#machine_" + machine_number)
                    .attr("id", "machine_" + machine_number + "_deadCards");
                $("#showIssues .deadCards").append(nav_item_deadCards);
                $("#machine_" + machine_number + "_deadCards").append(" " + card);
                this.deadCards_helper = machine_number;
                this.numberOfdeadCards++;
            }
            else {
                $("#machine_" + machine_number + "_deadCards").append(" " + card);
                this.numberOfdeadCards++;
            }
        }
    };
    MachineAPP.prototype.addMachine = function () {
        var _this = this;
        $("#addMachine").find(".buttonToAdd").on("click", function (event) {
            //let nameToAdd = event.currentTarget.parentElement;
            var element = event.currentTarget.parentElement;
            var addMachineName = element.id.replace('machineToAdd_', '');
            $.ajax({
                type: "POST",
                url: api_ip + "/client",
                header: ({ authentication: _this.cookie.access_token }),
                crossDomain: true,
                dataType: "json",
                contentType: "application/json",
                data: JSON.stringify({ Name: addMachineName })
            });
            element.remove();
            _this.getMiners();
        });
    };
    MachineAPP.prototype.menageMachine = function (event) {
        var _this = this;
        var miner = $(event.target);
        var machine_name = miner.parents(".miner").attr("data-id");
        if (miner.hasClass("delete") || miner.hasClass("hard-reset") || miner.hasClass("on_off"))
            $("#confirmOperation .input-group").addClass("d-none");
        else
            $("#confirmOperation .input-group").removeClass("d-none");
        if (miner.hasClass("on_off")) {
            $(".TurnOff, .TurnOn").removeClass("d-none");
            $(".confirm_operation").addClass("d-none");
        }
        else {
            $(".TurnOff, .TurnOn").addClass("d-none");
            $(".confirm_operation").removeClass("d-none");
        }
        $("#confirmOperation").find(".btn-primary").on("click", function () {
            var reason = $("#reason_of_management").val();
            if (miner.hasClass("reboot")) {
                $.ajax({
                    type: "POST",
                    url: api_ip + "/reboot",
                    header: ({ authentication: _this.cookie.access_token }),
                    crossDomain: true,
                    dataType: "json",
                    contentType: "application/json"
                });
            }
            if (miner.hasClass("restart")) {
                $.ajax({
                    type: "POST",
                    url: api_ip + "/restart",
                    header: ({ authentication: _this.cookie.access_token }),
                    crossDomain: true,
                    dataType: "json",
                    contentType: "application/json",
                    data: JSON.stringify({ MinerName: machine_name, reason: reason })
                });
            }
            if (miner.hasClass("delete")) {
                $.ajax({
                    type: "POST",
                    url: api_ip + "/client/delete",
                    header: ({ authentication: _this.cookie.access_token }),
                    crossDomain: true,
                    dataType: "json",
                    contentType: "application/json",
                    data: JSON.stringify({ Name: machine_name })
                });
            }
            if (miner.hasClass("hard-reset")) {
                $.ajax({
                    type: "POST",
                    url: api_ip + "/arduino/add_reset",
                    header: ({ authentication: _this.cookie.access_token }),
                    crossDomain: true,
                    dataType: "json",
                    contentType: "application/json",
                    data: JSON.stringify({ MachineName: machine_name, Function: "reset" })
                });
            }
            $("#confirmOperation").modal("hide");
            $("#reason_of_management").val("");
            $("#confirmOperation").find(".btn-primary").unbind("click");
        });
        $("#confirmOperation").find(".TurnOff").on("click", function () {
            $.ajax({
                type: "POST",
                url: api_ip + "/arduino/add_reset",
                header: ({ authentication: _this.cookie.access_token }),
                crossDomain: true,
                dataType: "json",
                contentType: "application/json",
                data: JSON.stringify({ MachineName: machine_name, Function: "shutdown" })
            });
        });
        $("#confirmOperation").find(".TurnOn").on("click", function () {
            $.ajax({
                type: "POST",
                url: api_ip + "/arduino/add_reset",
                header: ({ authentication: _this.cookie.access_token }),
                crossDomain: true,
                dataType: "json",
                contentType: "application/json",
                data: JSON.stringify({ MachineName: machine_name, Function: "poweron" })
            });
        });
    };
    MachineAPP.prototype.statistics = function () {
        var data = this.temperatureData;
        $(".t1").text(data.last.t1);
        $(".t2").text(data.last.t2);
        $(".v1").text(data.last.h1);
        $(".v2").text(data.last.h2);
        $(".last-time").text(moment(data.last.t_formated).format('MMM D HH:mm'));
        this.refreshGraph(data.data);
    }; //Functions belonging to graph [START]
    MachineAPP.prototype.refreshTemperatureDate = function () {
        var _this = this;
        var api = api_ip.split(":");
        $.get(api[0] + ":" + api[1] + "/api.php", function (data) {
            _this.temperatureData = data;
            _this.temperatureData.last.l += 10000;
            //this.refreshResetTime();
        });
    };
    MachineAPP.prototype.refreshGraph = function (data) {
        var temp = data.map(function (t) { return { "x": new Date(parseInt(t.t + "000")), "y": t.t1 }; });
        var temp2 = data.map(function (t) { return { "x": new Date(parseInt(t.t + "000")), "y": t.t2 }; });
        new Chartist.Line('.chart-cs', {
            series: [
                {
                    name: 'Sensor 1',
                    data: temp
                }, {
                    name: 'Sensor 2',
                    data: temp2
                }]
        }, {
            axisX: {
                type: Chartist.FixedScaleAxis,
                divisor: 6,
                labelInterpolationFnc: function (value) {
                    return moment(value).format('MMM D HH') + ":00";
                }
            },
            plugins: [Chartist.plugins.tooltip({
                    transformTooltipTextFnc: function (a) {
                        var values = a.split(",");
                        return values[1] + "&#186C</br><small>" + moment(parseInt(values[0])).format('HH:mm') + "</small>";
                    }
                })]
        });
    };
    MachineAPP.prototype.loadPinsSettings = function (success) {
        var _this = this;
        $.get(api_ip + "/arduino/pins", function (data) {
            data = data.filter(function (el) {
                return el.Function != null;
            });
            _this.pinsSetting = data;
            if (typeof success == "function")
                success(data);
        });
    };
    MachineAPP.prototype.editHardware = function () {
        var _this = this;
        $("#failEditHardware").text("");
        this.addedPins = [];
        $.ajax({
            url: api_ip + "/settings",
            type: "get",
            dataType: 'json'
        }).done(function (data) {
            _this.pinsCount = parseInt(data.pinCount);
            $(".numberOfPins").val(data.pinCount);
            _this.loadPinsSettings(function () {
                _this.renderMachinesResetOptions();
                $("#pinsHolder").droppable({
                    drop: function (event, ui) {
                        var i = parseInt(ui.draggable[0].innerText);
                        if (_this.pins.indexOf(i) != -1)
                            return;
                        _this.pins.push(i);
                        var index = _this.pinsSetting.findIndex(function (el) {
                            return el.ID == i;
                        });
                        _this.pinsSetting.splice(index, 1);
                        _this.revert = false;
                        $(event.target).append(ui.draggable[0]);
                        $(ui.draggable[0]).css({ "top": "initial", "left": "initial" });
                        _this.renderPins();
                        $.ajax({
                            type: "POST",
                            url: api_ip + "/arduino/pins",
                            header: ({ authentication: _this.cookie.access_token }),
                            crossDomain: true,
                            dataType: "json",
                            contentType: "application/json",
                            data: JSON.stringify({ ID: i.toString(), MinerName: null, Function: null })
                        });
                    }
                });
                // Load number of pins
                _this.generatePins();
            });
        }).fail(function () {
            console.log("error when try to generate pins");
        });
    };
    MachineAPP.prototype.renderPins = function () {
        var _this = this;
        var holder = $("#pinsHolder").html("");
        this.pins.sort(function (a, b) { return a - b; });
        for (var i = 0; i < this.pins.length; i++) {
            holder.append($("<span>").addClass("draggable bg-warning d-inline-block").text(this.pins[i]));
        }
        $(".draggable").draggable({
            containment: ".ui-widget-content",
            revert: function () { return _this.revert; },
            helper: "clone",
            start: function () {
                _this.revert = true;
            }
        });
    };
    MachineAPP.prototype.renderMachinesResetOptions = function () {
        var _this = this;
        var singleContainer = $(".minerPinsHolder");
        $(".minerPinsHolder").text("");
        singleContainer.html("");
        var template = $("#machinesResetOptionsTemplate>div");
        var _loop_1 = function(i) {
            var reset = this_1.pinsSetting.find(function (el) {
                return el.MinerName == _this.minersNamesArray[i] && el.Function == "0" && parseInt(el.ID) <= _this.pinsCount;
            });
            var power = this_1.pinsSetting.find(function (el) {
                return el.MinerName == _this.minersNamesArray[i] && el.Function == "1" && parseInt(el.ID) <= _this.pinsCount;
            });
            var el = $(template.clone());
            el.attr("data-miner", this_1.minersNamesArray[i]);
            el.find("label").text(this_1.minersNamesArray[i]);
            if (typeof reset != "undefined")
                el.find(".reset").append($("<span style='position: relative; left: calc(50% - 17px); top: 2px;'>")
                    .addClass("draggable bg-warning d-inline-block")
                    .text(reset.ID));
            if (typeof power != "undefined")
                el.find(".power").append($("<span style='position: relative; left: calc(50% - 17px); top: 2px;'>")
                    .addClass("draggable bg-warning d-inline-block")
                    .text(power.ID));
            singleContainer.append(el);
        };
        var this_1 = this;
        for (var i = 0; i < this.minersNamesArray.length; i++) {
            _loop_1(i);
        }
        $(".minerPinsHolder .droppable").droppable({
            drop: function (event, ui) {
                var i = ui.draggable[0].innerText;
                var m = event.target.parentNode.dataset.miner;
                var func = $(event.target).hasClass("reset") ? "0" : "1";
                var index = _this.pinsSetting.findIndex(function (el) {
                    return el.MinerName == m && el.Function == func;
                });
                if (index != -1) {
                    _this.revert = true;
                    return;
                }
                var index2 = _this.pinsSetting.findIndex(function (el) {
                    return el.ID == i;
                });
                if (index2 != -1)
                    _this.pinsSetting.splice(index2, 1);
                _this.addedPins.push(i);
                var newData = { ID: i, MinerName: m, Function: func };
                _this.pinsSetting.push(newData);
                $.ajax({
                    type: "POST",
                    url: api_ip + "/arduino/pins",
                    header: ({ authentication: _this.cookie.access_token }),
                    crossDomain: true,
                    dataType: "json",
                    contentType: "application/json",
                    data: JSON.stringify(newData)
                });
                _this.pins.splice(_this.pins.findIndex(function (el) {
                    return el == parseInt(i);
                }), 1);
                $(event.target).append(ui.draggable[0]);
                _this.revert = false;
                _this.enableHardButons();
                $(ui.draggable[0]).css({ "top": 2, "left": "calc(50% - 17px)", "position": "relative" });
            }
        });
    };
    MachineAPP.prototype.generatePins = function () {
        this.pins = [];
        var _loop_2 = function(i) {
            if (this_2.pinsSetting.findIndex(function (el) {
                return el.ID == i.toString() && el.Function != null;
            }) == -1)
                this_2.pins.push(i);
            if (typeof this_2.pinsSetting[i] != "undefined" &&
                this_2.pinsSetting[i].MinerName != null) {
                this_2.addedPins.push(this_2.pinsSetting[i].ID);
            }
        };
        var this_2 = this;
        for (var i = 1; i <= this.pinsCount; i++) {
            _loop_2(i);
        }
        this.renderPins();
    };
    MachineAPP.prototype.changePinsNumber = function () {
        var _this = this;
        var value = $(".numberOfPins").val();
        var highest_value = 0;
        for (var _i = 0, _a = this.pinsSetting; _i < _a.length; _i++) {
            var i = _a[_i];
            if (i.ID > highest_value)
                highest_value = i.ID;
        }
        if (highest_value > value) {
            $("#failEditHardware").text("A pin with a higher value has been already assigned");
        }
        else {
            $("#failEditHardware").text("");
            $.ajax({
                type: "POST",
                url: api_ip + "/settings",
                header: ({ authentication: this.cookie.access_token }),
                crossDomain: true,
                dataType: "json",
                contentType: "application/json",
                data: JSON.stringify({ Name: "pinCount", Value: value })
            }).done(function () {
                _this.pinsCount = parseInt(value);
                _this.generatePins();
                _this.renderMachinesResetOptions();
            });
        }
    };
    MachineAPP.prototype.fillWalletsList = function () {
        var _this = this;
        $.ajax({
            url: api_ip + "/wallet",
            type: "get",
            dataType: 'json'
        }).done(function (respond) {
            _this.walletsList = [];
            for (var _i = 0, respond_1 = respond; _i < respond_1.length; _i++) {
                var i = respond_1[_i];
                _this.walletsList.push({
                    "ID": i.ID,
                    "WalletName": i.WalletName,
                    "Address": i.Address,
                    "Currency": i.Currency,
                    "IsDefault": i.IsDefault
                });
            }
        }).fail(function () {
            console.log("Fail when try connect with wallets list");
        });
    };
    MachineAPP.prototype.fill_wallet_menagment = function () {
        $("#wallets_menagment_container").text("");
        for (var _i = 0, _a = this.walletsList; _i < _a.length; _i++) {
            var i = _a[_i];
            var singleWallet = $("#walletManagmentTemplate > div").clone();
            singleWallet.attr("id", "wallet_" + i.WalletName);
            singleWallet.attr("data-id", i.ID);
            $("#wallets_menagment_container").append(singleWallet);
            $("#wallet_" + i.WalletName).find(".wallet_name").text(i.WalletName);
            $("#wallet_" + i.WalletName).find(".wallet_address").text(i.Address);
            $("#wallet_" + i.WalletName).find(".wallet_currency").text(i.Currency);
            if ($("#wallet_" + i.WalletName).find(".wallet_currency").text() == "ETH") {
                $("#wallet_" + i.WalletName).find(".wallet_isDefault").attr("name", "wallet_radio_eth");
            }
            if ($("#wallet_" + i.WalletName).find(".wallet_currency").text() == "ETC") {
                $("#wallet_" + i.WalletName).find(".wallet_isDefault").attr("name", "wallet_radio_etc");
            }
            if (i.IsDefault === "1")
                $("#wallet_" + i.WalletName).find(".wallet_isDefault").attr("checked", true);
        }
    };
    MachineAPP.prototype.addWallet = function () {
        $("#addEditEditWallet").find(".modal-title").text("Add wallet");
        $("#addEditWallet").attr("data-wallet-id", "-1");
        if (!$("#addEditWalletErrors").hasClass("d-none"))
            $(".addEditWalletErrors").addClass("d-none");
        $("#addEdit_wallet_name").val("");
        $("#addEdit_wallet_address").val("");
    };
    MachineAPP.prototype.editWallet = function (event) {
        $("#addEditWallet").find(".modal-title").text("Edit wallet");
        var wallet = $(event.target).parents(".row");
        $("#addEditWallet").attr("data-wallet-id", wallet.attr("data-id"));
        if (!$("#addEditWalletErrors").hasClass("d-none"))
            $(".addEditWalletErrors").addClass("d-none");
        $("#addEdit_wallet_name").val($("#" + wallet.attr("id")).find(".wallet_name").text());
        $("#addEdit_wallet_address").val($("#" + wallet.attr("id")).find(".wallet_address").text());
        $("#addEdit_wallet_currency").val($("#" + wallet.attr("id")).find(".wallet_currency").text());
    };
    MachineAPP.prototype.walletAction = function () {
        var _this = this;
        if ($.trim($("#addEdit_wallet_address").val()) == '' || $.trim($("#addEdit_wallet_name").val()) == '') {
            $(".addEditWalletErrors").removeClass("d-none");
        }
        else {
            var id_1 = $("#addEditWallet").attr("data-wallet-id");
            var index_1 = this.walletsList.findIndex(function (el) { return el.ID == id_1; });
            var data_1 = {
                "WalletName": $("#addEdit_wallet_name").val(),
                "Address": $("#addEdit_wallet_address").val(),
                "Currency": $("#addEdit_wallet_currency").val(),
                "IsDefault": (index_1 == -1) ? 0 : this.walletsList[index_1].IsDefault,
                "ID": parseInt(id_1)
            };
            if (index_1 == -1)
                delete data_1.ID;
            $.ajax({
                type: "POST",
                url: api_ip + "/wallet",
                header: ({ authentication: this.cookie.access_token }),
                crossDomain: true,
                dataType: "json",
                contentType: "application/json",
                data: JSON.stringify(data_1)
            }).done(function (respond) {
                if (index_1 == -1) {
                    data_1.ID = parseInt(respond);
                    _this.walletsList.push(data_1);
                }
                else {
                    _this.walletsList[index_1] = data_1;
                }
                _this.fill_wallet_menagment();
                $("#addEditWallet").modal("hide");
            });
        }
    };
    MachineAPP.prototype.deleteWallet = function (event) {
        var _this = this;
        var wallet = $(event.target).parents(".row").attr("id");
        $(".wallet_error").text("");
        $("#confirmOperation .input-group").addClass("d-none");
        $("#confirmOperation .confirm_operation").on("click", function () {
            if ($("#" + wallet).find(".wallet_isDefault").is(':checked')) {
                $(".wallet_error").text("The default wallet can not be removed");
            }
            else {
                var id_2 = parseInt($("#addEditWallet").attr("data-wallet-id"));
                var index = _this.walletsList.findIndex(function (el) { return el.ID == id_2; });
                _this.walletsList.splice(index, 1);
                $.ajax({
                    type: "POST",
                    url: api_ip + "/wallet/delete",
                    header: ({ authentication: _this.cookie.access_token }),
                    crossDomain: true,
                    dataType: "json",
                    contentType: "application/json",
                    data: JSON.stringify({ "ID": id_2 })
                }).done(function () {
                    _this.fill_wallet_menagment();
                    $("#confirmOperation").modal("hide");
                    $("#confirmOperation .confirm_operation").unbind("click");
                });
            }
        });
    };
    MachineAPP.prototype.isDefaultWallet = function (event) {
        var _this = this;
        var wallet = $(event.target).parents(".row").attr("id");
        var id = $(event.target).parents(".row").attr("data-id");
        var index = this.walletsList.findIndex(function (el) { return el.ID == id; });
        var data = {
            "Currency": $("#" + wallet).find(".wallet_currency").text(),
            "ID": parseInt($("#" + wallet).attr("data-id"))
        };
        $.ajax({
            type: "POST",
            url: api_ip + "/wallet/set_default",
            header: ({ authentication: this.cookie.access_token }),
            crossDomain: true,
            dataType: "json",
            contentType: "application/json",
            data: JSON.stringify(data)
        }).done(function () {
            for (var i = 0; i < _this.walletsList.length; i++) {
                if (data.Currency == _this.walletsList[i].Currency) {
                    _this.walletsList[i].IsDefault = "0";
                }
            }
            _this.walletsList[index].IsDefault = "1";
        });
    };
    MachineAPP.prototype.fillConfigList = function () {
        var _this = this;
        $.ajax({
            url: api_ip + "/claymore/config/basic",
            type: "get",
            dataType: 'json'
        }).done(function (respond) {
            _this.configsList = [];
            for (var _i = 0, respond_2 = respond; _i < respond_2.length; _i++) {
                var i = respond_2[_i];
                _this.configsList.push({
                    "ID": i.id,
                    "ConfigName": i.name,
                    "Params": i.params,
                    "Currency": i.currency
                });
            }
        }).fail(function () {
            console.log("Fail when try connect with wallets list");
        });
    };
    MachineAPP.prototype.fill_config_menagment = function () {
        $("#configs_menagment_container").text("");
        for (var _i = 0, _a = this.configsList; _i < _a.length; _i++) {
            var i = _a[_i];
            var singleConfig = $("#configManagmentTemplate > div").clone();
            singleConfig.attr("id", "config_" + i.ConfigName);
            singleConfig.attr("data-id", i.ID);
            $("#configs_menagment_container").append(singleConfig);
            $("#config_" + i.ConfigName).find(".config_name").text(i.ConfigName);
            $("#config_" + i.ConfigName).find(".config_params").text(i.Params);
            $("#config_" + i.ConfigName).find(".config_currency").text(i.Currency);
            if ($("#config_" + i.ConfigName).find(".config_currency").text() == "ETH") {
                $("#config_" + i.ConfigName).find(".config_isDefault").attr("name", "config_radio_eth");
            }
            if ($("#config_" + i.ConfigName).find(".config_currency").text() == "ETC") {
                $("#config_" + i.ConfigName).find(".config_isDefault").attr("name", "config_radio_etc");
            }
        }
    };
    MachineAPP.prototype.addConfig = function () {
        $("#addEditEditConfig").find(".modal-title").text("Add config");
        $("#addEditConfig").attr("data-config-id", "-1");
        if (!$("#addEditConfigErrors").hasClass("d-none"))
            $(".addEditConfigErrors").addClass("d-none");
        $("#addEdit_config_name").val("");
        $("#addEdit_config_params").val("");
    };
    MachineAPP.prototype.editConfig = function (event) {
        $("#addEditConfig").find(".modal-title").text("Edit config");
        var config = $(event.target).parents(".row");
        $("#addEditConfig").attr("data-config-id", config.attr("data-id"));
        if (!$("#addEditConfigErrors").hasClass("d-none"))
            $(".addEditConfigErrors").addClass("d-none");
        $("#addEdit_config_name").val($("#" + config.attr("id")).find(".config_name").text());
        $("#addEdit_config_params").val($("#" + config.attr("id")).find(".config_params").text());
        $("#addEdit_config_currency").val($("#" + config.attr("id")).find(".config_currency").text());
    };
    MachineAPP.prototype.configAction = function () {
        var _this = this;
        if ($.trim($("#addEdit_config_params").val()) == '' || $.trim($("#addEdit_config_name").val()) == '') {
            $(".addEditConfigErrors").removeClass("d-none");
        }
        else {
            var id_3 = $("#addEditConfig").attr("data-config-id");
            var data_2 = {
                "ID": parseInt(id_3),
                "name": $("#addEdit_config_name").val(),
                "Params": $("#addEdit_config_params").val(),
                "Currency": $("#addEdit_config_currency").val()
            };
            var index_2 = this.configsList.findIndex(function (el) { return el.ID == id_3; });
            if (index_2 == -1)
                delete data_2.ID;
            $.ajax({
                type: "POST",
                url: api_ip + "/claymore/config/basic",
                header: ({ authentication: this.cookie.access_token }),
                crossDomain: true,
                dataType: "json",
                contentType: "application/json",
                data: JSON.stringify(data_2)
            }).done(function (respond) {
                if (index_2 == -1) {
                    data_2.ID = parseInt(respond);
                    _this.configsList.push(data_2);
                }
                else {
                    var data_3 = {
                        "ID": parseInt(id_3),
                        "ConfigName": $("#addEdit_config_name").val(),
                        "Params": $("#addEdit_config_params").val(),
                        "Currency": $("#addEdit_config_currency").val()
                    };
                    _this.configsList[index_2] = data_3;
                }
                _this.fill_config_menagment();
                $("#addEditConfig").modal("hide");
            });
        }
    };
    MachineAPP.prototype.deleteConfig = function () {
        var _this = this;
        $(".config_error").text("");
        $("#confirmOperation .input-group").addClass("d-none");
        $("#confirmOperation .confirm_operation").on("click", function () {
            var id = parseInt($("#addEditConfig").attr("data-config-id"));
            var index = _this.configsList.findIndex(function (el) { return el.ID == id; });
            _this.configsList.splice(index, 1);
            $.ajax({
                type: "POST",
                url: api_ip + "/claymore/config/basic/delete",
                header: ({ authentication: _this.cookie.access_token }),
                crossDomain: true,
                dataType: "json",
                contentType: "application/json",
                data: JSON.stringify({ "ID": id })
            });
            _this.fill_config_menagment();
            $("#confirmOperation").modal("hide");
            $("#confirmOperation .confirm_operation").unbind("click");
        });
    };
    MachineAPP.prototype.changeWallet = function () {
        $("#menage_machine").find(".change_config").text("");
        $("#menage_machine").find(".change_wallet").text("");
        for (var _i = 0, _a = this.configsList; _i < _a.length; _i++) {
            var i = _a[_i];
            $("#menage_machine").find(".change_config").append("<option>" + i.ConfigName + "</option>");
        }
        for (var _b = 0, _c = this.walletsList; _b < _c.length; _b++) {
            var i = _c[_b];
            $("#menage_machine").find(".change_wallet").append("<option>" + i.WalletName + "</option>");
        }
    };
    MachineAPP._instance = new MachineAPP();
    return MachineAPP;
}());
