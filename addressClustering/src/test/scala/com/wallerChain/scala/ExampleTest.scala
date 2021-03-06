///*
// * Copyright (c) 2016-2019 "Neo4j Sweden, AB" [https://neo4j.com]
// *
// * Licensed under the Apache License, Version 2.0 (the "License");
// * you may not use this file except in compliance with the License.
// * You may obtain a copy of the License at
// *
// *     http://www.apache.org/licenses/LICENSE-2.0
// *
// * Unless required by applicable law or agreed to in writing, software
// * distributed under the License is distributed on an "AS IS" BASIS,
// * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// * See the License for the specific language governing permissions and
// * limitations under the License.
// *
// * Attribution Notice under the terms of the Apache License 2.0
// *
// * This work was created by the collective efforts of the openCypher community.
// * Without limiting the terms of Section 6, any Derivative Work that is not
// * approved by the public consensus process of the openCypher Implementers Group
// * should not be described as “Cypher” (and Cypher® is a registered trademark of
// * Neo4j Inc.) or as "openCypher". Extensions by implementers or prototypes or
// * proposals for change that have been documented or implemented should only be
// * described as "implementation extensions to Cypher" or as "proposed changes to
// * Cypher that are not yet approved by the openCypher community".
// */
//package org.opencypher.spark.examples
//
//import java.io.{ByteArrayOutputStream, PrintStream}
//import java.net.URI
//
//import org.junit.runner.RunWith
//import org.opencypher.okapi.testing.Bag._
//import org.scalatestplus.junit.JUnitRunner
//import org.scalatest.{BeforeAndAfterAll, FunSpec, Matchers}
//
//import scala.io.Source
//
//@RunWith(classOf[JUnitRunner])
//abstract class ExampleTest extends FunSpec with Matchers with BeforeAndAfterAll {
//
//  private val oldStdOut = System.out
//
//  protected val emptyOutput: String = ""
//
//  protected def validate(app: => Unit, expectedOut: URI): Unit = {
//    validate(app, Source.fromFile(expectedOut).mkString)
//  }
//
//  protected def validateBag(app: => Unit, expectedOut: URI): Unit = {
//    val source = Source.fromFile(expectedOut)
//    val expectedLines = source.getLines().toList
//    val appLines = capture(app).split(System.lineSeparator())
//    withClue(s"${appLines.mkString("\n")} not equal to ${expectedLines.mkString("\n")}") {
//      appLines.toBag shouldEqual expectedLines.toBag
//    }
//  }
//
//  protected def validate(app: => Unit, expectedOut: String): Unit = {
//    capture(app) shouldEqual expectedOut
//  }
//
//  private def capture(app: => Unit): String = {
//    val charset = "UTF-8"
//    val outCapture = new ByteArrayOutputStream()
//    val printer = new PrintStream(outCapture, true, charset)
//    Console.withOut(printer)(app)
//    outCapture.toString(charset)
//  }
//
//  override protected def afterAll(): Unit = {
//    System.setOut(oldStdOut)
//    super.afterAll()
//  }
//}