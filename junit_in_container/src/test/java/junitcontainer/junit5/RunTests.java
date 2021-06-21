package junitcontainer.junit5;

import static org.junit.platform.engine.discovery.ClassNameFilter.includeClassNamePatterns;

import java.io.IOException;
import java.io.PrintWriter;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;
import junitcontainer.Constants;
import org.junit.platform.console.ConsoleLauncher;
import org.junit.platform.engine.DiscoverySelector;
import org.junit.platform.engine.discovery.DiscoverySelectors;
import org.junit.platform.launcher.Launcher;
import org.junit.platform.launcher.LauncherDiscoveryRequest;
import org.junit.platform.launcher.TestPlan;
import org.junit.platform.launcher.core.LauncherDiscoveryRequestBuilder;
import org.junit.platform.launcher.core.LauncherFactory;
import org.junit.platform.launcher.listeners.SummaryGeneratingListener;
import org.junit.platform.reporting.legacy.xml.LegacyXmlReportGeneratingListener;

public class RunTests {

  private final String[] testPackages;
  private final Path testingRoot;

  public RunTests(Path testingRoot, String... testPackages) {
    this.testingRoot = testingRoot;
    this.testPackages = testPackages;
  }

  public void runAll() throws IOException {

    Files.createDirectories(testingRoot);

    try (PrintWriter writer = new PrintWriter(System.out)) {

      SummaryGeneratingListener listener = new SummaryGeneratingListener();
      LegacyXmlReportGeneratingListener xmlListener =
          new LegacyXmlReportGeneratingListener(testingRoot, writer);

      List<? extends DiscoverySelector> selectors =
          Arrays.stream(testPackages)
              .map(DiscoverySelectors::selectPackage)
              .collect(Collectors.toList());

      LauncherDiscoveryRequest request =
          LauncherDiscoveryRequestBuilder.request()
              .selectors(selectors)
              .filters(includeClassNamePatterns(".*Test"))
              .build();
      Launcher launcher = LauncherFactory.create();
      TestPlan testPlan = launcher.discover(request);
      launcher.registerTestExecutionListeners(listener, xmlListener);
      launcher.execute(request);

      listener.getSummary().printTo(writer);
    }
  }

  public void runLauncher() {
    ConsoleLauncher.execute(
        System.out,
        System.err,
        "--select-package=junitcontainer.junit5container",
        "--reports-dir=/Users/humphrej/work/bazel-recipes");
  }

  public static void main(String[] args) throws Exception {

    String testingRoot =
        Optional.ofNullable(System.getenv("TESTING_ROOT")).orElse(Constants.DEFAULT_TESTING_ROOT);

    new RunTests(Paths.get(testingRoot), "junitcontainer/junit5").runAll();
  }
}

