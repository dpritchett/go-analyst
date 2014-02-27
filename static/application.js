var Analyst = {
  initialize: function() {
    Analyst.bindQueryRunner();
    console.log("Analyst loaded.");
  },

  columnTemplate: "<tr><% _.each(cols, function(col){ %><th><%= col %></th> <% }); %></tr>",

  rowTemplate:     "<tr><% _.each(row, function(cell){ %><td><%= cell %></td> <% }); %></tr>",

  bindQueryRunner: function() {
    $("#executeButton").on('click', Analyst.executeQuery);
  },

  executeQuery: function() {
    $.post("/sql-query", { query: $("#queryString").val() }, Analyst.renderResults)
  },

  resultSet: {},

  renderResults: function(res) {
    Analyst.resultSet = JSON.parse(res);
    Analyst.renderColumns();
    Analyst.renderRows();
  },

  columns: function() {
    return Analyst.resultSet.Rowset.Columns;
  },

  rows: function() {
    return Analyst.resultSet.Rowset.Rows;
  },

  renderColumns: function () {
    $tbl = $("#queryResults table");
    $tbl.empty();
    $tbl.append(_.template(Analyst.columnTemplate, { cols: Analyst.columns() }) );
  },

  renderRows: function() {
    $tbl = $("#queryResults table tbody");
    _.each(Analyst.rows(),
        function(row){
          $tbl.append(_.template(Analyst.rowTemplate, { row: row }) );
        });
  }
}

$(Analyst.initialize);
