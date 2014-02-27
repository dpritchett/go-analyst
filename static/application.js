var Analyst = {
  initialize: function() {
    Analyst.bindQueryRunner();
    console.log("Analyst loaded.");
  },

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
    $tbl = $("#queryResults table")
    $tbl.empty();
    $tbl.append("<tr></tr>");
    $.each(Analyst.columns(),
        function(ix, el){
          $("tr", $tbl).append("<th>" + el + "</th>");
          console.log(el)
        });
  },

  renderRows: function() {
    $.each(Analyst.rows(), function(ix, el){ console.log(el) });

    $tbl = $("#queryResults table tbody")
    $.each(Analyst.rows(),
        function(ix, row){
          $row = $tbl.append("<tr></tr>");
          $.each(row,
            function(ix, cell){
              $row.append("<td></td>").text(cell);
            });
        });
  }
}
  
$(Analyst.initialize);
