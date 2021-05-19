package com.wallerChain.scala

import org.opencypher.spark.api._
import org.opencypher.okapi.api.graph._
import org.opencypher.okapi.neo4j.io.Neo4jConfig
import java.net.URI
import org.opencypher.okapi.api.util.ZeppelinSupport._
//import org.opencypher.okapi.neo4j.io._
import org.neo4j.driver.v1.{Driver, Session, StatementResult}
import org.opencypher.spark.api.io.neo4j.sync.Neo4jGraphMerge
import org.opencypher.okapi.neo4j.io.MetaLabelSupport._
import org.opencypher.okapi.neo4j.io._
import org.apache.spark.sql.DataFrame



/**
 * @author ${user.name}
 */
object MultiInput {

//  def connectNeo4j(dataFixture: String, uri: String = "bolt://localhost:7687"): Neo4jContext = {
//
//  }

  def main(args: Array[String]): Unit = {
    // Create CAPS caps
    implicit val caps: CAPSSession = CAPSSession.local() //create(spark)

    // Change pw back to blob
    val neo4jConfig = Neo4jConfig(URI.create("bolt://localhost:7687"), user = "neo4j", password = Some("password"), encrypted =  false)

    //val neo4jNamespace = Namespace("txGraph")
    //val neo4jSource = GraphSources.cypher.neo4j(neo4j.config)

    // Access the graphs via their qualified graph names
    //val txGraph = caps.catalog.graph("txGraph.graph")
    //caps.catalog.register(Namespace("txGraph"), GraphSources.cypher.neo4j(neo4jConfig))
    //caps.catalog.register(Namespace("txGraph"), GraphSources.cypher.neo4j(neo4jConfig))
    caps.registerSource(Namespace("txGraph"), GraphSources.cypher.neo4j(neo4jConfig))

    //val txGraph = caps.catalog.graph("txGraph.graph")

    // Print all available graphs
    println("*** Graph names in catalog ***")
    println(caps.catalog.graphNames.mkString(System.lineSeparator()))

    /**
      * Returns a query that creates a graph containing persons that live in the same city and
      * know each other via 1 to 2 hops. The created graph contains a CLOSE_TO relationship between
      * each such pair of persons and is stored in the caps catalog using the given graph name.
      */
    def multiInputQuery(fromGraph: String): String =
      s"""FROM GRAPH $fromGraph
         |MATCH (a1:Address)-[:SENDS]->(t)<-[:SENDS]-(a2:Address)
         |CONSTRUCT
         |  CREATE (a1)-[:IS_SAME]->(a2)
         |RETURN GRAPH
      """.stripMargin

    // Create the Address Graph by applying the Multi-input heuristic
    val addressGraph = caps.cypher(multiInputQuery(s"txGraph.$entireGraphName")).graph

    caps.catalog.store("addressGraph", addressGraph)

//    // Copy products graph from File-based PGDS to Neo4j PGDS
//    caps.cypher(
//      s"""
//         |CATALOG CREATE GRAPH Neo4j.addressGraph {
//         |  FROM GRAPH $namespace.$graphName RETURN GRAPH
//         |}
//     """.stripMargin)

    // Test graph access
    // 5) Execute Cypher query and print results
    val result = caps.cypher(
      s"""|FROM GRAPH txGraph.graph
         |MATCH (t:Transaction)
         |RETURN t
         """.stripMargin)

//    val some = result.toString

    // 6) Collect results into string by selecting a specific column.
    //    This operation may be very expensive as it materializes results locally.
    //val names: Set[String] = result.records.table.df.collect().map(_.getAs[String]("n_name")).toSet

    //println(names)

    // Query for multi-inputs and create new edges between addresses that belong together
    //  val addressGraph = txGraph.cypher(
    //  val addressGraph = caps.cypher(
    //    s"""|FROM GRAPH txGraph.$entireGraphName
    //       |MATCH (a1:Address)-[:SENDS]->(t)<-[:SENDS]-(a2:Address)
    //       |CONSTRUCT
    //       |  CREATE (a1)-[:IS_SAME]->(a2)
    //       |RETURN GRAPH
    //    """.stripMargin).graph
    //
    //  def foo(x : Array[String]) = x.foldLeft("")((a,b) => a + b)
    //
    //  def main(args : Array[String]) {
    //    println( "Hello World!" )
    //    println("concat arguments = " + foo(args))
    //  }
  }
}
