package junitcontainer.junit5;

import static com.google.common.truth.Truth.assertThat;

import junitcontainer.UnderTest;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

public class MagicNumberTest {

  @Test
  @DisplayName("it's magic")
  public void shouldBeMagic() {
    assertThat(UnderTest.magicNumber()).isEqualTo(-136982450);
  }

  @Test
  @DisplayName("it's magic 2")
  public void shouldBeMagic2() {
    assertThat(UnderTest.magicNumber()).isEqualTo(-136982450);
  }
}

