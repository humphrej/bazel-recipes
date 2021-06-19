package junitcontainer.junit4;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Optional;
import junitcontainer.Constants;
import org.apache.tools.ant.taskdefs.optional.junit.XMLJUnitResultFormatter;
import org.junit.internal.TextListener;
import org.junit.runner.Description;
import org.junit.runner.JUnitCore;

public class RunTests {

  private final Path testingRoot;
  private final Class<?>[] testClazz;

  public RunTests(Path testingRoot, Class<?>... testClazz) {
    this.testingRoot = testingRoot;
    this.testClazz = testClazz;
  }

  public void runAll() throws IOException {

    JUnitCore junit = new JUnitCore();
    junit.addListener(new TextListener(System.out));

    junit.addListener(
        new JUnitResultFormatterAsRunListener(new XMLJUnitResultFormatter()) {
          @Override
          public void testStarted(Description description) throws Exception {
            String reportDir = testingRoot.toString();
            formatter.setOutput(
                new FileOutputStream(
                    new File(reportDir, "TEST-" + description.getDisplayName() + ".xml")));
            super.testStarted(description);
          }
        });

    junit.run(testClazz);
  }

  public static void main(String[] args) throws Exception {
    String testingRoot =
        Optional.ofNullable(System.getenv("TESTING_ROOT")).orElse(Constants.DEFAULT_TESTING_ROOT);

    new RunTests(Paths.get(testingRoot), TestSuite.class).runAll();
  }
}

