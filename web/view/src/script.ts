var api_ip: string = "http://0.0.0.0:8099";
var api_login: string = "http://0.0.0.0:9940";

declare var $: any, moment: any, Chartist: any, Cookies:any;

interface Cookie{
  access_token: string;
  expires_in: string;
  token_type: string;
  refresh_token: string;
}

class MachineAPP {

  private static _instance: MachineAPP = new MachineAPP();

  private cookie: Cookie;
  private access_token: string;
  private minersArray: any;
  private minersNamesArray: string[] = [];
  private all_mhs: number = 0;
  private all_temp: number = 0;
  private machine_quantity: number = 0;
  private numberOfDeadMachines: number = 0;
  private overhiting_helper: string;
  private numberOfOverhiting: number = 0;
  private deadCards_helper: string;
  private numberOfdeadCards: number = 0;
  private listOfdeadCards: string = "";
  private temperatureData: any;
  private sortType: string = "name";
  private sortAliveData: any[] = [];
  private sortDeadData: any[] = [];
  private pins: any = [];
  private addedPins: string[] = [];
  private pinsCount: number;
  private pinsSetting: any = [];
  private revert: boolean = true;
  private changeWatch: boolean = false;
  private walletsList: any = [];
  private configsList: any = [];



  constructor() {
    if (MachineAPP._instance) {
      throw new Error("Error: Instantiation failed: Use MachineAPP.getInstance() instead of new.");
    }
    MachineAPP._instance = this;
    this.minersArray = [];
    $(document).ready(() => this.initHandlers());
  }

  private initHandlers() {
    //this.checkingIsEmpty();
    this.isLogged();
    this.fillConfigList();
    this.fillWalletsList();
    this.getMiners();
    $(document).on("click", "#navbarCollapse [data-target='#addMachine']", () => this.addMinersList());
    $(document).on("click", "#logInButton", () => this.logIn());
    $(document).on("click", ".reboot", (event: any) => this.menageMachine(event));
    $(document).on("click", ".restart", (event: any) => this.menageMachine(event));
    $(document).on("click", ".delete", (event: any) => this.menageMachine(event));
    $(document).on("click", ".on_off", (event: any) => this.menageMachine(event));
    $(document).on("click", ".hard-reset", (event: any) => this.menageMachine(event));
    $(document).on("click", ".editHardwareButton", () => this.editHardware());
    $(document).on("click", ".wallet_menagment", () => this.fill_wallet_menagment());
    $(document).on("click", ".config_menagment", () => this.fill_config_menagment());
    $(document).on("click", "#changePinsNumber", () => this.changePinsNumber());
    $(document).on("click", ".wallet_edit", (event: any) => this.editWallet(event));
    $(document).on("click", ".wallet_add", () => this.addWallet());
    $(document).on("click", ".wallet_delete", (event: any) => this.deleteWallet(event));
    $(document).on("click", ".wallet_isDefault", (event: any) => this.isDefaultWallet(event));
    $(document).on("click", ".config_edit", (event: any) => this.editConfig(event));
    $(document).on("click", ".config_add", () => this.addConfig());
    $(document).on("click", ".config_delete", () => this.deleteConfig());
    $("#addEditWallet").find(".btn-outline-success").on("click", () => this.walletAction());
    $("#addEditConfig").find(".btn-outline-success").on("click", () => this.configAction());
    $(document).on("click", ".menage_machine", () => this.changeWallet());


    //  $(document).on("click", ".edit", (event: any) => this.editMachine(event));
    //  $(document).on("click", ".machine_operation .toConsol", (event: any) => this.toConsole(event));

    $('#statistics').on('shown.bs.modal', () => this.statistics());

    this.sortHandlers();
    $(".filter").on("keyup", () => this.filter());

    this.loadPinsSettings(() => this.enableHardButons());
    $("#editHardwareModal").on("click", "#clearPins", () => {
      this.pinsSetting = [];
      this.generatePins();
      $.ajax({
        type: "POST",
        url: api_ip + "/arduino/pins/clear",
        header: ({authentication: this.cookie.access_token}),
        crossDomain: true,
        dataType: "json",
        contentType: "application/json",
      });
      this.renderMachinesResetOptions();
      this.renderPins();

    });

    this.refreshTemperatureDate();
    setInterval(() => this.refreshResetTime(), 1000);
  }

  public static getInstance(): MachineAPP {
    return MachineAPP._instance;
  }

  private isLogged(){
    this.cookie = JSON.parse(Cookies.get('authentication'));
    $.ajax({
      url: api_login+"/rest/logged",
      type: "get",
      contentType: "application/json",
      Authorization: "Bearer " + this.cookie.access_token,
      dataType: 'json'
    }).done((respond: any) => {
      console.log("respond: ",respond);
    });
  }

  private logIn(){
    let nick = $("#logInNick").val();
    let password = $("#logInPassword").val();
    //let remember = $("#logInRemember").prop("checked");
    $.ajax({
      type: "POST",
      url:  api_login + "/oauth2/token",
      crossDomain: true,
      dataType: "json",
      contentType: "application/json",
      data: JSON.stringify({
        grant_type:"password",
        username:nick,
        password:password,
        client_id:"testclient",
        client_secret:"testpass"
       })
    }).done((respond: any) => {
      Cookies.set('authentication', JSON.stringify(respond),
      {expiry: respond.expires_in});
      this.cookie = JSON.parse(Cookies.get('authentication'));
    });
  }

  private filter() {
    let written_machine = $(".filter:focus").val();
    $('.filter').each(function() {
      $('.filter').val(written_machine);
    })

    if (written_machine != "") {
      for (let name of this.minersNamesArray) {
        if (name.indexOf(written_machine) == -1) {
          $("#machine_" + name).parent().hide();
        }
        else $("#machine_" + name).parent().show();
      }
    }

    else {
      for (let name of this.minersNamesArray) {
        $("#machine_" + name).parent().show();
      }
    }
  }

  private sortHandlers() {

    $(".sortByName").on("click", () => { //Sort by name
      $(".sortByDropdown .active").removeClass("active");
      $(".sortByName").addClass("active");
      this.sortType = "name";
      this.sortBy(this.sortType);
    });

    $(".sortByRunning").on("click", () => { //Sort by RunningTime
      $(".sortByDropdown .active").removeClass("active");
      $(".sortByRunning").addClass("active");
      this.sortType = "runTime";
      this.sortBy(this.sortType);

    });
    $(".sortByMh").on("click", () => { //Sort by Mh/s
      $(".sortByDropdown .active").removeClass("active");
      $(".sortByMh").addClass("active");
      this.sortType = "mh";
      this.sortBy(this.sortType);

    });
    $(".sortByTemp").on("click", () => { //Sort by temperature
      $(".sortByDropdown .active").removeClass("active");
      $(".sortByTemp").addClass("active");
      this.sortType = "heat";
      this.sortBy(this.sortType);
    });
  }

  private sortBy(type: string) {

    let asc = true;
    if (type == "runTime" || type == "mh" || type == "heat") {
      asc = false;
    }

    this.sortAliveData.sort((a, b) => {
      if (a[type] < b[type])
        return asc ? -1 : 1;
      if (a[type] > b[type])
        return asc ? 1 : -1;
      return 0;
    });

    this.sortDeadData.sort((a, b) => {
      if (a[name] < b[name])
        return asc ? -1 : 1;
      if (a[name] > b[name])
        return asc ? 1 : -1;
      return 0;
    });

    let container = $("#machine_container");
    let deadContainer = $("#dead_machine_container");

    for (let m of this.sortAliveData) {
      container.append($("#machine_" + m.name).parent());
    }

    for (let m of this.sortDeadData) {
      deadContainer.append($("#machine_" + m.name).parent());
    }
  }

  private addMinersList() {
    $.ajax({
      url: api_ip + "/clients/disconected",
      type: "get",
      dataType: 'json'
    }).done((minersToAdd: any) => {
      $("#addMachine").removeClass("d-none");
      $("#addMachine .modal-body").text(" ");


      let addMachineContainer = $(".addMachineList").clone();

      for (let i of minersToAdd) {
        if (i != "") {
          addMachineContainer.find(".machinesToAdd").first().attr("id", "machineToAdd_" + i);
          addMachineContainer.find(".nameToAdd").text(i);
          addMachineContainer.find(".machinesToAdd").first().removeClass("d-none").addClass("d-flex");
          $("#addMachine .modal-body").append(addMachineContainer.html());
        }
      }
      this.addMachine();
      if (minersToAdd[0] == "" || minersToAdd.length == 0) {
        $("#addMachine .modal-body").text("No machines to add");
      }

    }).fail(function() {
      $("#addMachine").addClass("d-none")
    });

  }

  private refreshResetTime() {
    if (typeof this.temperatureData == "undefined") return;

    // setTimeout(() => {
    //   this.refreshMinersData();
    // }, 5000)
  }

  private getMiners(): void { //Downloading miners
    $.ajax({
      url: api_ip + "/clients",
      type: "get",
      dataType: 'json'
    }).done((miners: any) => {
      if (!$("#alert-infos").hasClass("d-none")) {
        $("#alert-infos").addClass("d-none");
      }
      $("#dead_machine_container").text("");
      $("#dead_machine_container").addClass("d-none");
      $("#machine_container").text("");

      this.minersArray = miners;
      this.renderMiners();

      setInterval(() => {
        this.refreshMinersData();
      }, 12000);
      this.refreshMinersData();
    }).fail(function() {
      $("#alert-infos").removeClass("d-none");
      $("#alert-infos").text("Error when try connect with machines");
    });

  }

  private renderMiners(): void {
    for (let miner_key of Object.keys(this.minersArray)) { //Numer każdej koparki
      var machine: any = $('#template_machine>').clone();
      machine.find(">div").addClass("miner");
      machine.find(">div").first().attr("id", "machine_" + miner_key);
      this.minersNamesArray.push(miner_key);
      machine.find(">div").first().attr("data-id", miner_key);
      machine.find(".miner_name").text(this.minersArray[miner_key].MinerName);
      $("#machine_container").append(machine);

    };

    //Function needed to show popups correctly
    $(".editHardwareButton").addClass("d-md-flex");

  }

  private refreshMinersData(): void { //Function is clearing all data which need to be clear and fill it again
    $.ajax({
      url: api_ip + "/latest/stat",
      type: "get",
      dataType: 'json'
    }).done((respond: any) => {
      this.all_mhs = 0;
      this.all_temp = 0;
      this.machine_quantity = 0;
      this.numberOfDeadMachines = 0  //Refresh errors [END]


      $("#showIssues .overhiting").text(""); //Refresh errors [START]
      $("#showIssues .deadCards").text("");
      $(".issuesButton").addClass("d-none");

      $("#overhiting_counter").text("");
      $("#deadCards_counter").text("");
      $("#overhiting_span").addClass("d-none");
      $("#deadCards_span").addClass("d-none");
      $("#deadMachines_span").addClass("d-none");
      this.overhiting_helper = " ";
      this.deadCards_helper = " ";
      this.listOfdeadCards = " ";
      this.numberOfOverhiting = 0;
      this.numberOfdeadCards = 0;
      this.sortAliveData = [];
      this.sortDeadData = [];

      for (let miner_key of Object.keys(this.minersArray)) {
        let miner = this.minersArray[miner_key];
        let gpus = respond.hardware[miner_key];
        gpus = typeof
        gpus == "undefined" ? gpus = null : gpus = JSON.parse(respond.hardware[miner_key].gpuList);
        if (respond.miners[miner.MinerName] === undefined || gpus===null) { //When get data is fail
          this.deadMiners("machine_" + miner_key);
          continue;
        }
        this.fillMiners("machine_" + miner_key, respond.miners[miner.MinerName],gpus ,respond.config[miner_key]);
      }
      if (this.changeWatch == true) {
        this.sortBy(this.sortType);
      }

      this.all_mhs = this.Round(this.all_mhs, 0); //Refreshing info of total Mh/s and Temp on menu [START]
      this.all_temp = this.Round(this.all_temp / this.machine_quantity, 0);

      $(document).find(".total_data").html("<span class='d-block d-md-inline-block'>Total:  " + this.all_mhs + "Mh/s | "
        + "<span data-tooltip='yes' title='Averge temperature from all machines'>" + this.all_temp + "&#186C | </span></span>"
        + "<span data-tooltip='yes' title='Working machines'> Machines: <img src='style/icons/icon_machine_works.svg' alt='Working machines' width='20' height='15'/> "
        + this.machine_quantity + "</span>"
        + "<span data-tooltip='yes' title='Not working machines' > <img src='style/icons/icon_machine_not_works.svg' alt='Not working machines' width='20' height='15'/> "
        + this.numberOfDeadMachines) + "</span>"; //[END]

      if ((this.numberOfOverhiting + this.numberOfdeadCards + this.numberOfDeadMachines) > 0) {
        //If no errors, in modal we can see comunicate 'No errors'
        $(".issuesButton").removeClass("d-none");
        $(".numberOfIssues").text("Issues(" + (this.numberOfOverhiting + this.numberOfdeadCards + this.numberOfDeadMachines) + ")");
      }

      $("[data-tooltip='yes']").tooltip({ boundary: 'window' });
      this.changeWatch = false;
    }).fail(() => {
      console.log("Error when try refresh");
    });
  }

  public Round(n: number, k: number): number {
    var factor: number = Math.pow(10, k);
    return Math.round(n * factor) / factor;
  }

  public measureTime(minutes: number): any {
    let day = moment.duration(minutes, "minutes").humanize();
    let hours = moment.duration(minutes / 60, "minutes").humanize();
    let seconds = moment.duration(minutes / 120, "minutes").humanize();
    let day_test = /day/i;

    if (day_test.test(day)) {
      return day + " and " + hours;
    }
    else {
      return day + " and " + seconds;
    }

  }

  private fillMiners(id: string, miner: any, fans: any, wallet_info: any): void {
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
    if(wallet_info.Currency == "ETH"){
      $("#" + id).find(".currency").html('<img src="style/icons/icon_eth.svg" alt="ETH Graphic"/>' + wallet_info.Currency)
    }
    if(wallet_info.Currency == "ETC"){
      $("#" + id).find(".currency").html('<img src="style/icons/icon_etc.svg" alt="ETC Graphic"/>' + wallet_info.Currency)
    }
    $("#" + id).find(".wallet").text(" " + wallet_info.WalletName)
    let machine_name = id.replace("machine_", "");
    $("#" + id).find(".miner_name").text(machine_name);

    let arrayMghs: Array<any> = miner.DetailedEthHashRatePerGPU.split(';');
    let arrayTemp: Array<any> = miner.Temperatures.split(';');
    let forTemp = 0;
    let avergeTemp: any = 0;
    let oneGpu: any;
    let maxTemp = 0;

    for (let i: number = 0; i < arrayMghs.length; i++) {

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
        .val(fans[i + 1].declaredFanSpeed).change((event: any) => this.fanSpeedChanged(event));

      oneGpu.find(".fan_speed").addClass("fan_speed_" + machine_name + "_" + i);
      oneGpu.find(".fan_speed_" + machine_name + "_" + i).text(fans[i + 1].fanSpeed);

      forTemp += 2;
    }
    $("#" + id).find(".collapseButton").click((event: any) => {
      if (!$(event.currentTarget).parents(".mobile").find(".h-0").hasClass("show")) {
        $(event.currentTarget).parents(".gpus_container").find(".collapseArrow").css('transform', 'rotate(' + 360 + 'deg)');;
        $(event.currentTarget).parents(".gpus_container").find(".show").removeClass("show");
        $(event.currentTarget).parents(".mobile").find(".h-0").addClass("show");
        $(event.currentTarget).parents(".mobile").find(".collapseArrow").css('transform', 'rotate(' + 180 + 'deg)');
      }
      else {
        $(event.currentTarget).parents(".mobile").find(".h-0").removeClass("show");
        $(event.currentTarget).parents(".mobile").find(".collapseArrow").css('transform', 'rotate(' + 360 + 'deg)');;
      }

    });


    avergeTemp = this.Round((avergeTemp / arrayMghs.length), 0);
    let totalMhs: number = this.Round((miner.TotalEthHashRate / 1024), 1);
    let totalShares: number = miner.EthShares;
    let runningTime: number = this.measureTime(miner.RunningTime)

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
  }


  private enableHardButons() {
    $(".hard-reset, .on_off").addClass("d-none");
    for (let i = 0; i < this.pinsSetting.length; i++) {
      if (this.pinsSetting[i].Function != null) {
        $("#machine_" + this.pinsSetting[i].MinerName + " ." + ((this.pinsSetting[i].Function == 0) ? "hard-reset" : "on_off")).removeClass("d-none");
      }
    }
  }

  private fanSpeedChanged(event: any) {
    $.ajax({
      type: "POST",
      url: api_ip + "/fans",
      header: ({authentication: this.cookie.access_token}),
      crossDomain: true,
      dataType: "json",
      contentType: "application/json",
      data: JSON.stringify({
        id: event.target.dataset.gpu,
        machine: event.target.dataset.miner,
        speed: event.target.value
      })
    });
  }

  private deadMiners(id: string): void {

    $("#deadMachines_span").removeClass("d-none");
    $("#showIssues .deadMachines").removeClass("d-none");

    let deadMinersTemplete: any = $('#template_dead_machine').clone().html();
    let machine_name = id.replace("machine_", "");
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
  }

  private checkingErrors(oneGpu: any, arrayTemp: number, id: string, card: string, card_value: any) {
    $(oneGpu).find(".card-status").removeClass("bg-primary bg-dark bg-success bg-warning bg-danger text-dark");
    $(oneGpu).find(".row").removeClass("alert-dark text-dark");
    $(oneGpu).find(".gpu_info").text("");
    let machine_number = id.replace("machine_", "");

    if (arrayTemp == 0) {
      $(oneGpu).find(".card-status").addClass("bg-secondary");
      this.deadCards_stat(oneGpu, machine_number, card)
    }
    if (arrayTemp > 0 && arrayTemp <= 50) $(oneGpu).find(".card-status").addClass("bg-primary");
    if (arrayTemp > 50 && arrayTemp <= 65) $(oneGpu).find(".card-status").addClass("bg-success");
    if (arrayTemp > 65 && arrayTemp <= 75) $(oneGpu).find(".card-status").addClass("bg-warning");
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
      this.deadCards_stat(oneGpu, machine_number, card)
    }
  }

  private overhiting_stat(oneGpu: any, machine_number: string, card: string) {
    $("#overhiting_span").removeClass("d-none");
    $("#showIssues .overhiting").removeClass("d-none");
    $(oneGpu).find(".gpu_info").html("<img src='style/icons/icon_fire.svg' alt='Fire!' width='15px' data-tooltip='yes' title='Card is burning!'/>");

    if (this.overhiting_helper != machine_number) { //Avoid clone same machine name
      let nav_item_overhiting = $("<div>")
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

  }

  private deadCards_stat(oneGpu: any, machine_number: string, card: string) {
    $(oneGpu).find(".row").removeClass("alert-primary alert-dark alert-success alert-warning alert-danger bg-danger bg-secondary text-dark");
    $(oneGpu).find(".row").addClass("alert-dark text-dark");
    $(oneGpu).find(".gpu_info").html("<img src='style/icons/icon_dead.svg' alt='Dead card!' width='20px' data-tooltip='yes' title='Card disabled'/>");
    $("#showIssues #deadCards_span").removeClass("d-none");
    $("#showIssues .deadCards").removeClass("d-none");


    if (!$(document).find(".deadCards").hasClass("d-none")) {

      if (this.deadCards_helper != machine_number) { //Avoid clone same machine name
        let nav_item_deadCards = $("<div>")
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

  }

  public addMachine() {
    $("#addMachine").find(".buttonToAdd").on("click", (event: any) => {
      //let nameToAdd = event.currentTarget.parentElement;
      let element = event.currentTarget.parentElement;
      let addMachineName = element.id.replace('machineToAdd_', '');

      $.ajax({
        type: "POST",
        url: api_ip + "/client",
        header: ({authentication: this.cookie.access_token}),
        crossDomain: true,
        dataType: "json",
        contentType: "application/json",
        data: JSON.stringify({ Name: addMachineName })
      });

      element.remove();
      this.getMiners();
    });
  }

  public menageMachine(event: any) {

    let miner = $(event.target);
    let machine_name: string = miner.parents(".miner").attr("data-id");

    if (miner.hasClass("delete") || miner.hasClass("hard-reset") || miner.hasClass("on_off")) $("#confirmOperation .input-group").addClass("d-none");
    else $("#confirmOperation .input-group").removeClass("d-none");

    if (miner.hasClass("on_off")) {
      $(".TurnOff, .TurnOn").removeClass("d-none");
      $(".confirm_operation").addClass("d-none");
    }
    else {
      $(".TurnOff, .TurnOn").addClass("d-none");
      $(".confirm_operation").removeClass("d-none");
    }

    $("#confirmOperation").find(".btn-primary").on("click", () => {
      let reason: string = $("#reason_of_management").val();
      if (miner.hasClass("reboot")) {
        $.ajax({
          type: "POST",
          url: api_ip + "/reboot",
          header: ({authentication: this.cookie.access_token}),
          crossDomain: true,
          dataType: "json",
          contentType: "application/json",
         });
      }

      if (miner.hasClass("restart")) {
        $.ajax({
          type: "POST",
          url: api_ip + "/restart",
          header: ({authentication: this.cookie.access_token}),
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
          header: ({authentication: this.cookie.access_token}),
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
          header: ({authentication: this.cookie.access_token}),
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

    $("#confirmOperation").find(".TurnOff").on("click", () => {
      $.ajax({
        type: "POST",
        url: api_ip + "/arduino/add_reset",
        header: ({authentication: this.cookie.access_token}),
        crossDomain: true,
        dataType: "json",
        contentType: "application/json",
        data: JSON.stringify({ MachineName: machine_name, Function: "shutdown" })
      });
    });


    $("#confirmOperation").find(".TurnOn").on("click", () => {
      $.ajax({
        type: "POST",
        url: api_ip + "/arduino/add_reset",
        header: ({authentication: this.cookie.access_token}),
        crossDomain: true,
        dataType: "json",
        contentType: "application/json",
        data: JSON.stringify({ MachineName: machine_name, Function: "poweron" })
      });
    });
  }

  public statistics() {
    let data = this.temperatureData;
    $(".t1").text(data.last.t1);
    $(".t2").text(data.last.t2);
    $(".v1").text(data.last.h1);
    $(".v2").text(data.last.h2);
    $(".last-time").text(moment(data.last.t_formated).format('MMM D HH:mm'));
    this.refreshGraph(data.data);
  }   //Functions belonging to graph [START]

  private refreshTemperatureDate() {
    let api = api_ip.split(":");
    $.get(api[0] + ":" + api[1] + "/api.php", (data: any) => {
      this.temperatureData = data;

      this.temperatureData.last.l += 10000;
      //this.refreshResetTime();
    });
  }


  private refreshGraph(data: any) {
    let temp = data.map((t: any) => { return { "x": new Date(parseInt(t.t + "000")), "y": t.t1, }; });
    let temp2 = data.map((t: any) => { return { "x": new Date(parseInt(t.t + "000")), "y": t.t2, }; });

    new Chartist.Line('.chart-cs', {
      series: [
        {
          name: 'Sensor 1',
          data: temp
        }, {
          name: 'Sensor 2',
          data: temp2
        }]
    }
      , {
        axisX: {
          type: Chartist.FixedScaleAxis,
          divisor: 6,
          labelInterpolationFnc: function(value: any) {
            return moment(value).format('MMM D HH') + ":00";
          }
        },
        plugins: [Chartist.plugins.tooltip({
          transformTooltipTextFnc: (a: any) => {
            let values = a.split(",");
            return values[1] + "&#186C</br><small>" + moment(parseInt(values[0])).format('HH:mm') + "</small>";
          }
        })],
      });

  }

  private loadPinsSettings(success: any) {
    $.get(api_ip + "/arduino/pins", (data: any) => {
      data = data.filter((el: any) => {
        return el.Function != null;
      });
      this.pinsSetting = data;
      if (typeof success == "function")
        success(data);
    });
  }

  private editHardware() {
    $("#failEditHardware").text("");
    this.addedPins = [];
    $.ajax({
      url: api_ip + "/settings",
      type: "get",
      dataType: 'json'
    }).done((data: any) => {
      this.pinsCount = parseInt(data.pinCount);
      $(".numberOfPins").val(data.pinCount);

      this.loadPinsSettings(() => {
        this.renderMachinesResetOptions();
        $("#pinsHolder").droppable({
          drop: (event: any, ui: any) => {
            let i = parseInt(ui.draggable[0].innerText);

            if (this.pins.indexOf(i) != -1) return;
            this.pins.push(i);
            let index = this.pinsSetting.findIndex((el: any) => {
              return el.ID == i;
            });
            this.pinsSetting.splice(index, 1);
            this.revert = false;
            $(event.target).append(ui.draggable[0]);
            $(ui.draggable[0]).css({ "top": "initial", "left": "initial" });
            this.renderPins();

            $.ajax({
              type: "POST",
              url: api_ip + "/arduino/pins",
              header: ({authentication: this.cookie.access_token}),
              crossDomain: true,
              dataType: "json",
              contentType: "application/json",
              data: JSON.stringify({ ID: i.toString(), MinerName: null, Function: null })
            });
          }
        });


        // Load number of pins
        this.generatePins();
      });
    }).fail(() => {
      console.log("error when try to generate pins");
    });
  }

  private renderPins() {

    let holder: any = $("#pinsHolder").html("");
    this.pins.sort((a: number, b: number) => { return a - b; });

    for (let i = 0; i < this.pins.length; i++) {
      holder.append(
        $("<span>").addClass("draggable bg-warning d-inline-block").text(this.pins[i])
      );
    }

    $(".draggable").draggable({
      containment: ".ui-widget-content",
      revert: () => { return this.revert; },
      helper: "clone",
      start: () => {
        this.revert = true;
      }
    });

  }

  private renderMachinesResetOptions() {
    let singleContainer = $(".minerPinsHolder");
    $(".minerPinsHolder").text("");
    singleContainer.html("");
    let template = $("#machinesResetOptionsTemplate>div");

    for (let i = 0; i < this.minersNamesArray.length; i++) {
      let reset = this.pinsSetting.find((el: any) => {
        return el.MinerName == this.minersNamesArray[i] && el.Function == "0" && parseInt(el.ID) <= this.pinsCount;
      });
      let power = this.pinsSetting.find((el: any) => {
        return el.MinerName == this.minersNamesArray[i] && el.Function == "1" && parseInt(el.ID) <= this.pinsCount;
      });

      let el = $(template.clone());
      el.attr("data-miner", this.minersNamesArray[i]);
      el.find("label").text(this.minersNamesArray[i]);

      if (typeof reset != "undefined")
        el.find(".reset").append($("<span style='position: relative; left: calc(50% - 17px); top: 2px;'>")
          .addClass("draggable bg-warning d-inline-block")
          .text(reset.ID));
      if (typeof power != "undefined")
        el.find(".power").append($("<span style='position: relative; left: calc(50% - 17px); top: 2px;'>")
          .addClass("draggable bg-warning d-inline-block")
          .text(power.ID));

      singleContainer.append(el);
    }

    $(".minerPinsHolder .droppable").droppable({
      drop: (event: any, ui: any) => {
        let i: string = ui.draggable[0].innerText;
        let m: string = event.target.parentNode.dataset.miner;
        let func: string = $(event.target).hasClass("reset") ? "0" : "1";
        let index = this.pinsSetting.findIndex((el: any) => {
          return el.MinerName == m && el.Function == func;
        });

        if (index != -1) {
          this.revert = true;
          return;
        }

        let index2 = this.pinsSetting.findIndex((el: any) => {
          return el.ID == i;
        });
        if (index2 != -1)
          this.pinsSetting.splice(index2, 1);

        this.addedPins.push(i);
        let newData = { ID: i, MinerName: m, Function: func };
        this.pinsSetting.push(newData);
        $.ajax({
          type: "POST",
          url: api_ip + "/arduino/pins",
          header: ({authentication: this.cookie.access_token}),
          crossDomain: true,
          dataType: "json",
          contentType: "application/json",
          data: JSON.stringify(newData)
        });



        this.pins.splice(this.pins.findIndex((el: number) => {
          return el == parseInt(i);
        }), 1);

        $(event.target).append(ui.draggable[0]);
        this.revert = false;
        this.enableHardButons();
        $(ui.draggable[0]).css({ "top": 2, "left": "calc(50% - 17px)", "position": "relative" });
      }
    });
  }

  private generatePins() {
    this.pins = [];
    for (let i = 1; i <= this.pinsCount; i++) {
      if (this.pinsSetting.findIndex((el: any) => {
        return el.ID == i.toString() && el.Function != null;
      }) == -1)
        this.pins.push(i);

      if (typeof this.pinsSetting[i] != "undefined" &&
        this.pinsSetting[i].MinerName != null) {
        this.addedPins.push(this.pinsSetting[i].ID);
      }
    }
    this.renderPins();
  }

  private changePinsNumber() {
    let value = $(".numberOfPins").val();
    let highest_value: number = 0;
    for (let i of this.pinsSetting) {
      if (i.ID > highest_value) highest_value = i.ID;
    }
    if (highest_value > value) {
      $("#failEditHardware").text("A pin with a higher value has been already assigned");
    }
    else {
      $("#failEditHardware").text("");
      $.ajax({
        type: "POST",
        url: api_ip + "/settings",
        header: ({authentication: this.cookie.access_token}),
        crossDomain: true,
        dataType: "json",
        contentType: "application/json",
        data: JSON.stringify({ Name: "pinCount", Value: value })
      }).done(() => {
        this.pinsCount = parseInt(value);
        this.generatePins()
        this.renderMachinesResetOptions();
      });
    }
  }
  private fillWalletsList() {
    $.ajax({
      url: api_ip + "/wallet",
      type: "get",
      dataType: 'json'
    }).done((respond: any) => {
      this.walletsList = [];
      for (let i of respond) {
        this.walletsList.push({
          "ID": i.ID,
          "WalletName": i.WalletName,
          "Address": i.Address,
          "Currency": i.Currency,
          "IsDefault": i.IsDefault
        });
      }
    }).fail(function() {
      console.log("Fail when try connect with wallets list")
    });
  }

  private fill_wallet_menagment() {
    $("#wallets_menagment_container").text("");
    for (let i of this.walletsList) {

      let singleWallet = $("#walletManagmentTemplate > div").clone();
      singleWallet.attr("id", "wallet_" + i.WalletName);
      singleWallet.attr("data-id", i.ID);
      $("#wallets_menagment_container").append(singleWallet);

      $("#wallet_" + i.WalletName).find(".wallet_name").text(i.WalletName);
      $("#wallet_" + i.WalletName).find(".wallet_address").text(i.Address);
      $("#wallet_" + i.WalletName).find(".wallet_currency").text(i.Currency);

      if ($("#wallet_" + i.WalletName).find(".wallet_currency").text() == "ETH") {
        $("#wallet_" + i.WalletName).find(".wallet_isDefault").attr("name", "wallet_radio_eth")
      }
      if ($("#wallet_" + i.WalletName).find(".wallet_currency").text() == "ETC") {
        $("#wallet_" + i.WalletName).find(".wallet_isDefault").attr("name", "wallet_radio_etc")
      }
      if (i.IsDefault === "1") $("#wallet_" + i.WalletName).find(".wallet_isDefault").attr("checked", true);

    }
  }

  private addWallet() {
    $("#addEditEditWallet").find(".modal-title").text("Add wallet");
    $("#addEditWallet").attr("data-wallet-id", "-1");
    if (!$("#addEditWalletErrors").hasClass("d-none")) $(".addEditWalletErrors").addClass("d-none");

    $("#addEdit_wallet_name").val("");
    $("#addEdit_wallet_address").val("");
  }

  private editWallet(event: any) {
    $("#addEditWallet").find(".modal-title").text("Edit wallet");
    let wallet = $(event.target).parents(".row");
    $("#addEditWallet").attr("data-wallet-id", wallet.attr("data-id"));

    if (!$("#addEditWalletErrors").hasClass("d-none")) $(".addEditWalletErrors").addClass("d-none");
    $("#addEdit_wallet_name").val($("#" + wallet.attr("id")).find(".wallet_name").text());
    $("#addEdit_wallet_address").val($("#" + wallet.attr("id")).find(".wallet_address").text());
    $("#addEdit_wallet_currency").val($("#" + wallet.attr("id")).find(".wallet_currency").text());
  }

  private walletAction() {
    if ($.trim($("#addEdit_wallet_address").val()) == '' || $.trim($("#addEdit_wallet_name").val()) == '') {
      $(".addEditWalletErrors").removeClass("d-none");
    }
    else {
      let id = $("#addEditWallet").attr("data-wallet-id");

      let index = this.walletsList.findIndex((el: any) => el.ID == id);

      let data = {
        "WalletName": $("#addEdit_wallet_name").val(),
        "Address": $("#addEdit_wallet_address").val(),
        "Currency": $("#addEdit_wallet_currency").val(),
        "IsDefault": (index == -1)? 0 : this.walletsList[index].IsDefault,
        "ID": parseInt(id)
      };
      if (index == -1) delete data.ID;

      $.ajax({
        type: "POST",
        url: api_ip + "/wallet",
        header: ({authentication: this.cookie.access_token}),
        crossDomain: true,
        dataType: "json",
        contentType: "application/json",
        data: JSON.stringify(data)
      }).done((respond: string) => {
        if (index == -1) { // Add
          data.ID = parseInt(respond);
          this.walletsList.push(data);
        } else { // Edit
          this.walletsList[index] = data;
        }
        this.fill_wallet_menagment();
        $("#addEditWallet").modal("hide");
      });
    }
  }


  private deleteWallet(event: any) {
    let wallet = $(event.target).parents(".row").attr("id");

    $(".wallet_error").text("");
    $("#confirmOperation .input-group").addClass("d-none")
    $("#confirmOperation .confirm_operation").on("click", () => {
      if ($("#" + wallet).find(".wallet_isDefault").is(':checked')) {
        $(".wallet_error").text("The default wallet can not be removed");
      }
      else {
        let id: number = parseInt($("#addEditWallet").attr("data-wallet-id"));

        let index = this.walletsList.findIndex((el: any) => el.ID == id);
        this.walletsList.splice(index, 1);
        $.ajax({
          type: "POST",
          url: api_ip + "/wallet/delete",
          header: ({authentication: this.cookie.access_token}),
          crossDomain: true,
          dataType: "json",
          contentType: "application/json",
          data: JSON.stringify({ "ID": id })
        }).done(() => {
          this.fill_wallet_menagment();
          $("#confirmOperation").modal("hide");
          $("#confirmOperation .confirm_operation").unbind("click");
        });
      }
    });
  }

  private isDefaultWallet(event: any) {
    let wallet = $(event.target).parents(".row").attr("id");
    let id = $(event.target).parents(".row").attr("data-id");
    let index = this.walletsList.findIndex((el: any) => el.ID == id);
    let data = {
      "Currency": $("#" + wallet).find(".wallet_currency").text(),
      "ID": parseInt($("#" + wallet).attr("data-id"))
    };

    $.ajax({
      type: "POST",
      url: api_ip + "/wallet/set_default",
      header: ({authentication: this.cookie.access_token}),
      crossDomain: true,
      dataType: "json",
      contentType: "application/json",
      data: JSON.stringify(data)
    }).done(() => {
      for (let i = 0; i < this.walletsList.length; i++) {
        if (data.Currency == this.walletsList[i].Currency) {
          this.walletsList[i].IsDefault = "0";
        }
      }
      this.walletsList[index].IsDefault = "1";
    })
  }

  private fillConfigList() {
    $.ajax({
      url: api_ip + "/claymore/config/basic",
      type: "get",
      dataType: 'json'
    }).done((respond: any) => {
      this.configsList = [];
      for (let i of respond) {
        this.configsList.push({
          "ID": i.id,
          "ConfigName": i.name,
          "Params": i.params,
          "Currency": i.currency
        });
      }
    }).fail(function() {
      console.log("Fail when try connect with wallets list")
    });
  }

  private fill_config_menagment() {
    $("#configs_menagment_container").text("");
    for (let i of this.configsList) {
      let singleConfig = $("#configManagmentTemplate > div").clone();
      singleConfig.attr("id", "config_" + i.ConfigName);
      singleConfig.attr("data-id", i.ID);
      $("#configs_menagment_container").append(singleConfig);

      $("#config_" + i.ConfigName).find(".config_name").text(i.ConfigName);
      $("#config_" + i.ConfigName).find(".config_params").text(i.Params);
      $("#config_" + i.ConfigName).find(".config_currency").text(i.Currency);

      if ($("#config_" + i.ConfigName).find(".config_currency").text() == "ETH") {
        $("#config_" + i.ConfigName).find(".config_isDefault").attr("name", "config_radio_eth")
      }
      if ($("#config_" + i.ConfigName).find(".config_currency").text() == "ETC") {
        $("#config_" + i.ConfigName).find(".config_isDefault").attr("name", "config_radio_etc")
      }
    }
  }

  private addConfig() {
    $("#addEditEditConfig").find(".modal-title").text("Add config");
    $("#addEditConfig").attr("data-config-id", "-1");
    if (!$("#addEditConfigErrors").hasClass("d-none")) $(".addEditConfigErrors").addClass("d-none");

    $("#addEdit_config_name").val("");
    $("#addEdit_config_params").val("");
  }

  private editConfig(event: any) {
    $("#addEditConfig").find(".modal-title").text("Edit config");
    let config = $(event.target).parents(".row");
    $("#addEditConfig").attr("data-config-id", config.attr("data-id"));

    if (!$("#addEditConfigErrors").hasClass("d-none")) $(".addEditConfigErrors").addClass("d-none");
    $("#addEdit_config_name").val($("#" + config.attr("id")).find(".config_name").text());
    $("#addEdit_config_params").val($("#" + config.attr("id")).find(".config_params").text());
    $("#addEdit_config_currency").val($("#" + config.attr("id")).find(".config_currency").text());
  }

  private configAction() {
    if ($.trim($("#addEdit_config_params").val()) == '' || $.trim($("#addEdit_config_name").val()) == '') {
      $(".addEditConfigErrors").removeClass("d-none");
    }
    else {
      let id = $("#addEditConfig").attr("data-config-id");
      let data = {
        "ID": parseInt(id),
        "name": $("#addEdit_config_name").val(),
        "Params": $("#addEdit_config_params").val(),
        "Currency": $("#addEdit_config_currency").val()
      };
      let index = this.configsList.findIndex((el: any) => el.ID == id);
      if (index == -1) delete data.ID;

      $.ajax({
        type: "POST",
        url: api_ip + "/claymore/config/basic",
        header: ({authentication: this.cookie.access_token}),
        crossDomain: true,
        dataType: "json",
        contentType: "application/json",
        data: JSON.stringify(data)
      }).done((respond: string) => {
        if (index == -1) { // Add
          data.ID = parseInt(respond);
          this.configsList.push(data);
        } else { // Edit
          let data = {
            "ID": parseInt(id),
            "ConfigName": $("#addEdit_config_name").val(),
            "Params": $("#addEdit_config_params").val(),
            "Currency": $("#addEdit_config_currency").val()
          }
          this.configsList[index] = data;
        }

        this.fill_config_menagment();
        $("#addEditConfig").modal("hide");
      });
    }
  }

  private deleteConfig() {

    $(".config_error").text("");
    $("#confirmOperation .input-group").addClass("d-none")
    $("#confirmOperation .confirm_operation").on("click", () => {
      let id: number = parseInt($("#addEditConfig").attr("data-config-id"));
      let index = this.configsList.findIndex((el: any) => el.ID == id);
      this.configsList.splice(index, 1);
      $.ajax({
        type: "POST",
        url: api_ip + "/claymore/config/basic/delete",
        header: ({authentication: this.cookie.access_token}),
        crossDomain: true,
        dataType: "json",
        contentType: "application/json",
        data: JSON.stringify({ "ID": id })
      });
      this.fill_config_menagment();
      $("#confirmOperation").modal("hide");
      $("#confirmOperation .confirm_operation").unbind("click");
    });
  }

  private changeWallet(){
    $("#menage_machine").find(".change_config").text("");
    $("#menage_machine").find(".change_wallet").text("");

     for(let i of this.configsList){
       $("#menage_machine").find(".change_config").append("<option>"+i.ConfigName+"</option>")
     }
     for(let i of this.walletsList){
       $("#menage_machine").find(".change_wallet").append("<option>"+i.WalletName+"</option>")
     }
  }
}
