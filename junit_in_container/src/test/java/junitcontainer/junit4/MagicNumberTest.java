package junitcontainer.junit4;

import static com.google.common.truth.Truth.assertThat;

import junitcontainer.UnderTest;
import org.junit.Test;

public class MagicNumberTest {

  @Test
  public void shouldBeMagic() {
    assertThat(UnderTest.magicNumber()).isEqualTo(-136982450);
  }
}

