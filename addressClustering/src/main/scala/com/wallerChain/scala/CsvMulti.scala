package com.wallerChain.scala

import org.opencypher.spark.api._
import org.opencypher.okapi.api.graph._
import org.apache.spark.sql.{DataFrame, functions}
import org.opencypher.okapi.api.graph.CypherResult
import org.opencypher.okapi.neo4j.io.Neo4jConfig
import java.net.URI

import org.opencypher.okapi.api.util.ZeppelinSupport._
import org.opencypher.okapi.impl.util.PrintOptions
import org.opencypher.spark.api.io.fs.FSGraphSource
import org.opencypher.okapi.neo4j.io.testing.Neo4jTestUtils._
import org.opencypher.spark.api.io.neo4j.Neo4jPropertyGraphDataSource
//import org.opencypher.okapi.neo4j.io._
import org.neo4j.driver.v1.{Driver, Session, StatementResult}
import org.opencypher.spark.api.io.neo4j.sync.Neo4jGraphMerge
import org.opencypher.okapi.neo4j.io.MetaLabelSupport._
import org.opencypher.okapi.neo4j.io._
import org.apache.spark.sql.DataFrame
import org.opencypher.okapi.impl.util


/**
  * @author ${user.name}
  */
object CsvMultiInput extends App {

  //  def connectNeo4j(dataFixture: String, uri: String = "bolt://localhost:7687"): Neo4jContext = {
  //
  //  }

  //def main(args: Array[String]): Unit = {
    // Create CAPS session
    implicit val session: CAPSSession = CAPSSession.local() //create(spark)

    val neo4jConfig = Neo4jConfig(URI.create("bolt://localhost:7687"), user = "neo4j", password = Some("password"), encrypted =  false)


    //val graphDir = getClass.getResource("/fs-graphsource/csv").getFile
    val graphDir = getClass.getResource("/fs-graphsource/caps").getFile

  println("Graph dir: " + graphDir)

    // Create File-based PGDS
    val filePgds: FSGraphSource = GraphSources.fs(rootPath = graphDir).csv
    val neo4jPgds: Neo4jPropertyGraphDataSource = GraphSources.cypher.neo4j(neo4jConfig)

    // Register PGDS in the catalog
    val namespace = Namespace("CSV")
    session.registerSource(namespace, filePgds)
    session.registerSource(Namespace("Neo4j"), neo4jPgds)

    // Print graphs stored in PGDS
    println("*** Graph names in File-based PGDS ***")
    println(filePgds.graphNames.mkString(System.lineSeparator()))

    // Print all available graphs
    println("*** Graph names in catalog ***")
    println(session.catalog.graphNames.mkString(System.lineSeparator()))

    // Get graph name from File-based PGDS
    val graphName = filePgds.graphNames.head

  println("Graph name: " + graphName)

  // Copy products graph from File-based PGDS to Neo4j PGDS
  session.cypher(
    s"""
       |CATALOG CREATE GRAPH Neo4j.$graphName {
       |  FROM GRAPH $namespace.$graphName RETURN GRAPH
       |}
     """.stripMargin)

    // Access graph via Cypher query
    //session.cypher(s"FROM GRAPH $namespace.$graphName MATCH (n) RETURN n LIMIT 20")//.show
//    val results: CypherResult = session.cypher(s"FROM GRAPH $namespace.$graphName MATCH (n) RETURN n LIMIT 20")

    // 4) Extract DataFrame representing the query result
    //val df: DataFrame = results.records.asDataFrame
    //results.records.show(PrintOptions.out)
  //results.show

    // Access graph via API
//    session.catalog.graph(QualifiedGraphName(namespace, graphName)).cypher("MATCH (n) RETURN n LIMIT 20")//.show

    //val neo4jConfig = Neo4jConfig(URI.create("bolt://localhost:7687"), user = "neo4j", password = Some("blob"), encrypted = false)

    //val neo4j = connectNeo4j()


  //}
}